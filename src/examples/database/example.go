package main

import (
	"github.com/berlincount/go-jsonstruct"

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
	person_struct_json := `
        {"struct": "person",
         "fields": [
          {"name": "FirstName","type": "string", "tags": "db:\"first_name\""},
          {"name": "LastName", "type": "string", "tags": "db:\"last_name\""}
        ]}
	`
	person_struct_type, _ := jsonstruct.Decode(strings.NewReader(person_struct_json))
	person_struct_value := reflect.New(person_struct_type)
	person_struct_interface := person_struct_value.Interface()

	// populate with some data
	insertPerson := `INSERT INTO person (first_name, last_name) VALUES (?, ?)`
	db.MustExec(insertPerson, "Elaine", "Marley")
	db.MustExec(insertPerson, "Fester", "Shinetop")
	db.MustExec(insertPerson, "Herman", "Toothrot")

	// query the data into Go struct
	rows, _ := db.Queryx("SELECT * FROM person")
	for rows.Next() {
		_ = rows.StructScan(person_struct_interface)
		spew.Dump(person_struct_interface)
	}
}
