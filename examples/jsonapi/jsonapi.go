package main

import (
	"github.com/berlincount/jsonstruct"

	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/jsonapi"
)

func main() {
	// JSONAPI uses an extended Datatype that we like
	jsonstruct.MapType("time.Time", reflect.TypeOf(time.Now()))

	// Remarshal JSONAPI data as created by testBlogForCreate in
	// https://github.com/google/jsonapi/blob/master/examples/app.go
	blogStructJSON := `
        {"struct": "comment",
         "fields": [
          {"name": "ID",     "type": "int",    "tags": "jsonapi:\"primary,comments\""},
          {"name": "PostID", "type": "int",    "tags": "jsonapi:\"attr,post_id\""},
          {"name": "Body",   "type": "string", "tags": "jsonapi:\"attr,body\""}
        ]}
        {"struct": "post",
         "fields": [
          {"name": "ID",       "type": "int", "tags": "jsonapi:\"primary,posts\""},
          {"name": "BlogID",   "type": "int", "tags": "jsonapi:\"attr,blog_id\""},
          {"name": "Title",    "type": "string", "tags": "jsonapi:\"attr,title\""},
          {"name": "Body",     "type": "string", "tags": "jsonapi:\"attr,body\""},
          {"name": "Comments", "type": "[]*comment", "tags": "jsonapi:\"relation,comments\""}
        ]}
        {"struct": "blog",
         "fields": [
         {"name": "ID", "type": "int", "tags": "jsonapi:\"primary,blogs\""},
         {"name": "Title", "type": "string", "tags": "jsonapi:\"attr,title\""},
         {"name": "Posts", "type": "[]*post", "tags": "jsonapi:\"relation,posts\""},
         {"name": "CurrentPost", "type": "*post", "tags": "jsonapi:\"relation,current_post\""},
         {"name": "CurrentPostID", "type": "int", "tags": "jsonapi:\"attr,current_post_id\""},
         {"name": "CreatedAt", "type": "time.Time", "tags": "jsonapi:\"attr,created_at\""},
         {"name": "ViewCount", "type": "int", "tags": "jsonapi:\"attr,view_count\""}
        ]}
	`
	decodedStructs, err := jsonstruct.Decode(strings.NewReader(blogStructJSON))
	if err != nil {
		panic("something went terribly wrong decoding the structures")
	}
	if len(decodedStructs) != 3 {
		panic("something went rather unexpected decoding the structures")
	}
	blogStructType := decodedStructs["blog"]
	blogStructValue := reflect.New(blogStructType)
	blogStructInterface := blogStructValue.Interface()

	err = jsonapi.UnmarshalPayload(strings.NewReader(`
	{"data":{"type":"blogs","id":"1","attributes":{"created_at":1480279183,"current_post_id":0,"title":"Title 1","view_count":0},"relationships":{"current_post":{"data":{"type":"posts","id":"1","attributes":{"blog_id":0,"body":"Bar","title":"Foo"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"2","attributes":{"body":"bar","post_id":0}}]}}}},"posts":{"data":[{"type":"posts","id":"1","attributes":{"blog_id":0,"body":"Bar","title":"Foo"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"2","attributes":{"body":"bar","post_id":0}}]}}},{"type":"posts","id":"2","attributes":{"blog_id":0,"body":"Bas","title":"Fuubar"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"3","attributes":{"body":"bas","post_id":0}}]}}}]}}}}
	`), blogStructInterface)
	if err != nil {
		panic("something went terribly wrong unmarshalling the payload")
	}

	// data can be found in orderly structures with same signature as in original example now

	buf := bytes.NewBuffer(nil)
	err = jsonapi.MarshalOnePayloadEmbedded(buf, blogStructInterface)
	if err != nil {
		panic("something went terribly wrong remarshalling the payload")
	}

	fmt.Println(buf.String())
}
