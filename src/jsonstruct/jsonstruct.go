package jsonstruct

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

// Type registry

var TypeMap map[string]reflect.Type

func MapType(Name string, Type reflect.Type) {
	// Register Type under Name
	TypeMap[Name] = Type
	// Register derived common types
	TypeMap["[]"+Name] = reflect.SliceOf(Type)
	TypeMap["*"+Name] = reflect.PtrTo(Type)
	TypeMap["[]*"+Name] = reflect.SliceOf(reflect.PtrTo(Type))
}

func init() {
	TypeMap = make(map[string]reflect.Type)
	// approximations for JSON Datatypes Number, String & Boolean
	MapType("int", reflect.TypeOf(0))
	MapType("float", reflect.TypeOf(.0))
	MapType("string", reflect.TypeOf(""))
	MapType("bool", reflect.TypeOf(true))
}

// Decoding

type Field struct {
	Name string            "json:name"
	Type string            "json:type"
	Tags reflect.StructTag "json:tags"
}
type Struct struct {
	Struct string "json:struct"
	Fields []Field
}

func Decode(r io.Reader) (reflect.Type, error) {
	dec := json.NewDecoder(r)

	// we're reconstructing a stream of one or more structs
	var m Struct
	var recon_struct reflect.Type = reflect.TypeOf(nil)
	for {
		// catch JSON decode errors
		if err := dec.Decode(&m); err == io.EOF {
			// EOF? return the last struct type we found
			return recon_struct, nil
		} else if err != nil {
			return nil, err
		}

		// JSON data inconsistent?
		if len(m.Struct) < 0 {
			return nil, errors.New("empty struct name")
		}

		// gather fields of struct
		new_struct := make([]reflect.StructField, 0, len(m.Fields))
		for _, field := range m.Fields {
			new_struct = append(new_struct, reflect.StructField{field.Name, "", TypeMap[field.Type], field.Tags, 0, nil, false})
		}

		// create new struct type (and register it)
		recon_struct = reflect.StructOf(new_struct)
		MapType(m.Struct, recon_struct)
	}
	return recon_struct, nil
}
