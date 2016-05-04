package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/confirm"
	_ "gopkg.in/authboss.v0/lock"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/register"
	_ "gopkg.in/authboss.v0/remember"

	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"

	"gopkg.in/pg.v4"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var (
	ab        = authboss.New()
	templates = tpl.Must(tpl.Load("views/app", "views/partials", "layout.html.tpl", funcs))
	schemaDec = schema.NewDecoder()
)

func configureAuthBoss(db *pg.DB, config *Config) {
	database := NewPostgresStorer(db)
	ab.Storer = database
	ab.MountPath = "/auth"
	ab.ViewsPath = "./views/auth"
	ab.RootURL = config.Web.Root
	ab.LayoutDataMaker = layoutData
	b, err := ioutil.ReadFile(filepath.Join("views/app", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b)))
	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}
	ab.CookieStoreMaker = NewCookieStorer
	ab.SessionStoreMaker = NewSessionStorer
	ab.Mailer = authboss.SMTPMailer(config.Mailer.Server, nil)
	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       8,
			AllowWhitespace: false,
		},
	}
	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	// wire stuff up
	config, err := LoadConfig("./conf/application.yml")
	if err != nil {
		panic(err)
	}

	db := pg.Connect(&pg.Options{
		User:     config.Database.User,
		Password: config.Database.Password,
		Database: config.Database.Database,
	})

	cookieStoreKey := []byte(config.Crypto.Application)
	sessionStoreKey := []byte(config.Crypto.Session)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)

	configureAuthBoss(db, config)

	schemaDec.IgnoreUnknownKeys(true)

	mux := mux.NewRouter()
	mux.PathPrefix("/auth").Handler(ab.NewRouter())
	mux.Handle("/secure", protect(secureArea)).Methods("GET")
	mux.HandleFunc("/", index).Methods("GET")

	mux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not found")
	})

	// Set up our middleware chain
	stack := alice.New(logger, nosurfing, ab.ExpireMiddleware).Then(mux)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

type Digit struct {
	Position int64
	Digit    int64
}

func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r)
	mustRender(w, r, "index", data)
}

func secureArea(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r)
	mustRender(w, r, "secure", data)
}

func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	data.MergeKV("csrf_token", nosurf.Token(r))
	err := templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}
