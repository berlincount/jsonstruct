// Package jsonstruct provides a JSON deserializer for Go structures for Go 1.7+
package jsonstruct

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
)

// TypeMap provides a Type registry, mapping type names to reflect Types
var TypeMap map[string]reflect.Type

// MapType adds a new Type (plus its [], * and []* variants) to the registry
func MapType(Name string, Type reflect.Type) {
	// Register Type under Name
	TypeMap[Name] = Type
	// Register derived common types
	TypeMap["[]"+Name] = reflect.SliceOf(Type)
	TypeMap["*"+Name] = reflect.PtrTo(Type)
	TypeMap["[]*"+Name] = reflect.SliceOf(reflect.PtrTo(Type))
}

// initalize basic types we can expect with JSON
func init() {
	TypeMap = make(map[string]reflect.Type)
	// approximations for JSON Datatypes Number, String & Boolean
	MapType("int", reflect.TypeOf(0))
	MapType("float", reflect.TypeOf(.0))
	MapType("string", reflect.TypeOf(""))
	MapType("bool", reflect.TypeOf(true))
}

// Field holds a JSON description of individual Go fields
type Field struct {
	Name string            "json:\"name\""
	Type string            "json:\"type\""
	Tags reflect.StructTag "json:\"tags\""
}

// Struct holds JSON description of Go structures
type Struct struct {
	Struct string "json:\"struct\""
	Fields []Field
}

// Decode one or multiple Go structures from JSON, register and return their Types
func Decode(r io.Reader) ([]reflect.Type, error) {
	dec := json.NewDecoder(r)

	// we're reconstructing a stream of one or more structs
	var m Struct
	var reconStruct []reflect.Type
	for {
		// catch JSON decode errors
		if err := dec.Decode(&m); err == io.EOF {
			// EOF? return our collected struct types
			return reconStruct, nil
		} else if err != nil {
			return nil, err
		}

		// JSON data inconsistent?
		if len(m.Struct) <= 0 {
			return nil, errors.New("empty struct name")
		}

		// gather fields of struct
		newStruct := make([]reflect.StructField, 0, len(m.Fields))
		for _, field := range m.Fields {
			firstChar := strings.Split(field.Name, "")[0]
			if firstChar == strings.ToLower(firstChar) {
				return nil, errors.New("unsupported private fields found in structures")
			}
			newStruct = append(newStruct, reflect.StructField{
				Name:      field.Name,
				PkgPath:   "",
				Type:      TypeMap[field.Type],
				Tag:       field.Tags,
				Offset:    0,
				Index:     nil,
				Anonymous: false})
		}

		// create new struct type (and register it)
		reconStruct = append(reconStruct, reflect.StructOf(newStruct))
		MapType(m.Struct, reconStruct[len(reconStruct)-1])

		// continue in loop until EOF (error condition) is encountered
	}
}
