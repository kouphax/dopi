package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/authboss.v0"
	"gopkg.in/pg.v4"
)

type User struct {
	ID   int
	Name string

	// Auth
	Email    string
	Password string

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Lock
	AttemptNumber int64
	AttemptTime   time.Time
	Locked        time.Time

	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time

	// Remember is in another table
}

type Token struct {
	Key   string
	Token string
}

type PostgresStorer struct {
	db *pg.DB
}

func NewPostgresStorer(db *pg.DB) *PostgresStorer {
	return &PostgresStorer{
		db: db,
	}
}

func (s PostgresStorer) Create(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	err := s.db.Create(&user)

	spew.Dump(user)
	return err
}

func (s PostgresStorer) Put(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	// TODO - need to derive the primary key

	return s.db.Update(&user)
}

func (s PostgresStorer) Get(key string) (result interface{}, err error) {
	var user User
	err = s.db.Model(&user).Where("email = ?", key).First()
	if err != nil {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s PostgresStorer) AddToken(key, token string) error {
	return s.db.Create(&Token{
		Key:   key,
		Token: token,
	})
}

func (s PostgresStorer) DelTokens(key string) error {
	var token Token
	_, err := s.db.Model(&token).Where("key = ?", key).Delete()
	return err
}

func (s PostgresStorer) UseToken(givenKey, token string) error {
	var t Token
	result, err := s.db.Model(&t).Where("key = ?", givenKey).Where("token = ?", token).Delete()
	if err != nil {
		return err
	}

	if result.Affected() == 0 {
		return authboss.ErrTokenNotFound
	}

	return nil
}

func (s PostgresStorer) ConfirmUser(tok string) (result interface{}, err error) {
	var user User
	err = s.db.Model(&user).Where("confirmToken = ?", tok).First()
	if err != nil {
		return nil, authboss.ErrUserNotFound
	}

	return user, nil
}

func (s PostgresStorer) RecoverUser(rec string) (result interface{}, err error) {
	var user User
	err = s.db.Model(&user).Where("recoverToken = ?", rec).First()
	if err != nil {
		return nil, authboss.ErrUserNotFound
	}

	return user, nil
}
