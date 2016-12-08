package main

import (
	"github.com/berlincount/jsonstruct"

	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// create a database (in-memory only)
	var db *sqlx.DB
	db = sqlx.MustConnect("sqlite3", ":memory:")

	// create database schema
	schema := `
	CREATE TABLE person (
	  first_name VARCHAR(20),
	  last_name  VARCHAR(20)
	);`
	db.MustExec(schema)

	// create Go struct from JSON
	/* equivalent to:
	type Person struct {
	  FirstName string `db:"first_name"`
	  LastName  string `db:"last_name"`
	}
	*/
	personStructJSON := `
        {"struct": "person",
         "fields": [
          {"name": "FirstName","type": "string", "tags": "db:\"first_name\""},
          {"name": "LastName", "type": "string", "tags": "db:\"last_name\""}
        ]}
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(personStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 1 {
		panic("something went rather unexpected decoding the structures")
	}
	personStructType := decodedStructs[len(decodedStructs)-1]
	personStructValue := reflect.New(personStructType)
	personStructInterface := personStructValue.Interface()

	// populate with some data
	insertPerson := `INSERT INTO person (first_name, last_name) VALUES (?, ?)`
	db.MustExec(insertPerson, "Elaine", "Marley")
	db.MustExec(insertPerson, "Fester", "Shinetop")
	db.MustExec(insertPerson, "Herman", "Toothrot")

	// query the data into Go struct
	rows, err := db.Queryx("SELECT * FROM person")
	if err != nil {
		panic("something went terribly wrong while selecting rows")
	}
	for rows.Next() {
		err = rows.StructScan(personStructInterface)
		if err != nil {
			panic("something went terribly wrong while fetching row")
		}
		spew.Dump(personStructInterface)
	}
}
