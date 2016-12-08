package jsonstruct_test

import (
	"github.com/berlincount/jsonstruct"

	"reflect"
	"strings"
	"testing"
)

type testStruct struct {
	myTestInt   int
	myTestField string
}

func TestMapType(t *testing.T) {
	jsonstruct.MapType("testStruct", reflect.TypeOf(testStruct{}))

	if len(jsonstruct.TypeMap)%4 != 0 {
		t.Errorf("jsonstruct.TypeMap expected to have four elements for each type")
	}

	if rtype, present := jsonstruct.TypeMap["testStruct"]; present {
		if rtype != reflect.TypeOf(testStruct{}) {
			t.Errorf("wrong basic type was registered")
		}
	} else {
		t.Errorf("basic type was not registered")
	}

	if rtype, present := jsonstruct.TypeMap["*testStruct"]; present {
		if rtype != reflect.TypeOf(new(testStruct)) {
			t.Errorf("wrong * type was registered")
		}
	} else {
		t.Errorf("* type was not registered")
	}

	if rtype, present := jsonstruct.TypeMap["[]testStruct"]; present {
		if rtype != reflect.TypeOf([]testStruct{}) {
			t.Errorf("wrong [] type was registered")
		}
	} else {
		t.Errorf("[] type was not registered")
	}

	if rtype, present := jsonstruct.TypeMap["[]*testStruct"]; present {
		if rtype != reflect.TypeOf([]*testStruct{}) {
			t.Errorf("wrong []* type was registered")
		}
	} else {
		t.Errorf("[]* type was not registered")
	}
}

func TestDecode1(t *testing.T) {
	testStructJSON := `
        {"struct": "test",
         "fields": [
          {"name": "TestInt",    "type": "int",    "tags": "testTag:\"first_field\""},
          {"name": "TestString", "type": "string", "tags": "testTag:\"second_field\""}
        ]}
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(testStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 1 {
		panic("something went rather unexpected decoding the structures")
	}

	if decodedStructs[0].NumField() != 2 {
		t.Errorf("wrong number of fields was decoded")
	}

	firstField := decodedStructs[0].Field(0)
	if firstField.Name != "TestInt" {
		t.Errorf("first field name was unpacked wrongly")
	}
	if firstField.Type != reflect.TypeOf(0) {
		t.Errorf("first field type was unpacked wrongly")
	}
	tag, ok := firstField.Tag.Lookup("testTag")
	if !ok {
		t.Errorf("first field tag was unpacked wrongly")
	}
	if tag != "first_field" {
		t.Errorf("first field tag was unpacked badly")
	}
	if firstField.Anonymous {
		t.Errorf("first field is anonymous")
	}

	secondField := decodedStructs[0].Field(1)
	if secondField.Name != "TestString" {
		t.Errorf("second field name was unpacked wrongly")
	}
	if secondField.Type != reflect.TypeOf("") {
		t.Errorf("second field type was unpacked wrongly")
	}
	tag, ok = secondField.Tag.Lookup("testTag")
	if !ok {
		t.Errorf("second field tag was unpacked wrongly")
	}
	if tag != "second_field" {
		t.Errorf("second field tag was unpacked badly")
	}
	if secondField.Anonymous {
		t.Errorf("second field is anonymous")
	}

	if rtype, present := jsonstruct.TypeMap["test"]; present {
		if rtype != decodedStructs[0] {
			t.Errorf("wrong basic type was registered")
		}
	} else {
		t.Errorf("basic type was not registered")
	}
}

func TestDecode2(t *testing.T) {
	testStructJSON := `
        {"struct": "test1",
         "fields": []}
        {"struct": "test2",
         "fields": []}
        {"struct": "test3",
         "fields": []}
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(testStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 3 {
		panic("something went rather unexpected decoding the structures")
	}

	if decodedStructs[0].NumField() != 0 || decodedStructs[1].NumField() != 0 || decodedStructs[2].NumField() != 0 {
		t.Errorf("wrong number of fields was decoded")
	}

	if rtype, present := jsonstruct.TypeMap["test1"]; present {
		if rtype != decodedStructs[0] {
			t.Errorf("wrong basic type1 was registered")
		}
	} else {
		t.Errorf("basic type1 was not registered")
	}

	if rtype, present := jsonstruct.TypeMap["test2"]; present {
		if rtype != decodedStructs[1] {
			t.Errorf("wrong basic type2 was registered")
		}
	} else {
		t.Errorf("basic type2 was not registered")
	}

	if rtype, present := jsonstruct.TypeMap["test3"]; present {
		if rtype != decodedStructs[2] {
			t.Errorf("wrong basic type3 was registered")
		}
	} else {
		t.Errorf("basic type3 was not registered")
	}
}

func TestDecode3(t *testing.T) {
	testStructJSON := `
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(testStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 0 {
		panic("something went rather unexpected decoding the structures")
	}
}

func TestDecode4(t *testing.T) {
	testStructJSON := `
        {"struct": "testP",
         "fields": [
          {"name": "private", "type": "int", "tags": ""}
        ]}
	`
	_, err := jsonstruct.Decode(strings.NewReader(testStructJSON))
	if err == nil {
		panic("not getting an error when trying to create struct with private field")
	}
}
