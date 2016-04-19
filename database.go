package main

import (
	"gopkg.in/pg.v4"
	"io/ioutil"
)

func readFile(path string) (error, string) {
	contentsAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err, ""
	}

	return nil, string(contentsAsBytes[:])
}

func create(db *pg.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err, contents := readFile("./db/carteblanche.sql")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = db.Exec(contents)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})

	err := create(db)
	if err != nil {
		panic(err)
	}
}
