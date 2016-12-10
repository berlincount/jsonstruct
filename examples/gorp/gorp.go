package main

import (
	"github.com/berlincount/jsonstruct"

	"encoding/json"
	"log"
	"reflect"
	"strings"

	"database/sql"
	// NOTE: we're using a fork that supports automatic index generation
	"github.com/kimxilxyong/gorp"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// initialize the DbMap
	dbmap, structmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// populate with some data (easier from JSON when using generated structures)
	personStructInterface := reflect.New(structmap["person"]).Interface()
	buf := []byte(`{"FirstName":"Elaine","LastName":"Marley"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	dbmap.Insert(personStructInterface)
	buf = []byte(`{"FirstName":"Fester","LastName":"Shinetop"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	dbmap.Insert(personStructInterface)
	buf = []byte(`{"FirstName":"Herman","LastName":"Toothrot"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	dbmap.Insert(personStructInterface)

	// use convenience SelectInt
	count, err := dbmap.SelectInt("select count(*) from persons")
	checkErr(err, "select count(*) failed")
	log.Println("Rows after inserting:", count)

	// fetch all rows
	persons := reflect.New(reflect.SliceOf(structmap["person"]))
	_, err = dbmap.Select(persons.Interface(), "select * from persons order by last_name")
	checkErr(err, "Select failed")
	log.Println("All rows:")
	for x := 0; x < persons.Elem().Len(); x++ {
		log.Printf("    %d: %v\n", x, persons.Elem().Index(x))
	}
}

func initDb() (*gorp.DbMap, map[string]reflect.Type) {
	db, err := sql.Open("sqlite3", ":memory:")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// create Go struct from JSON
	/* equivalent to:
	type Person struct {
	    FirstName string `db:"name:first_name,size:20"`
	    LastName  string `db:"name:last_name,size:20"`
	}
	*/
	personStructJSON := `
        {"struct": "person",
         "fields": [
          {"name": "FirstName","type": "string", "tags": "db:\"name:first_name,size:20\""},
          {"name": "LastName", "type": "string", "tags": "db:\"name:last_name,size:20\""}
        ]}
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(personStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 1 {
		panic("something went rather unexpected decoding the structures")
	}

	for table := range decodedStructs {
		// use plural struct name for table names
		dbmap.AddTableWithName(reflect.Zero(decodedStructs[table]).Interface(), table+"s")
	}

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap, decodedStructs
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
