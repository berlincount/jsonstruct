package main

import (
	"github.com/berlincount/jsonstruct"

	"encoding/json"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		// Madness!
		// We're extracting the table name from the dbtable tag of the
		// first field, as gorm can't find the struct name - and would
		// always use "" as table name
		tablename, present := reflect.TypeOf(db.Value).Elem().Field(0).Tag.Lookup("dbtable")
		if !present {
			return "unknown"
		}
		return tablename

	}
	// create a database (in-memory only)
	db, err := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// create Go struct from JSON
	/* equivalent to:
	type Person struct {
	  FirstName string `dbtable:"persons" db:"first_name"`
	  LastName  string `db:"last_name"`
	}
	*/
	personStructJSON := `
        {"struct": "person",
         "fields": [
          {"name": "FirstName","type": "string", "tags": "dbtable:\"persons\" db:\"first_name\""},
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

	// create or upgrade our DB schema automagically
	db.AutoMigrate(personStructInterface)

	// populate with some data (easier from JSON when using generated structures)
	buf := []byte(`{"FirstName":"Elaine","LastName":"Marley"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	db.Create(personStructInterface)
	buf = []byte(`{"FirstName":"Fester","LastName":"Shinetop"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	db.Create(personStructInterface)
	buf = []byte(`{"FirstName":"Herman","LastName":"Toothrot"}`)
	err = json.Unmarshal(buf, &personStructInterface)
	db.Create(personStructInterface)

	// query the data into Go struct
	db.First(personStructInterface)
	spew.Dump(personStructInterface)
}
