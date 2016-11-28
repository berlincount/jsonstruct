package main

import (
	"github.com/berlincount/go-jsonstruct"

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
	blog_struct_json := `
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
	blog_struct_type, _ := jsonstruct.Decode(strings.NewReader(blog_struct_json))
	blog_struct_value := reflect.New(blog_struct_type)
	blog_struct_interface := blog_struct_value.Interface()

	jsonapi.UnmarshalPayload(strings.NewReader(`
	{"data":{"type":"blogs","id":"1","attributes":{"created_at":1480279183,"current_post_id":0,"title":"Title 1","view_count":0},"relationships":{"current_post":{"data":{"type":"posts","id":"1","attributes":{"blog_id":0,"body":"Bar","title":"Foo"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"2","attributes":{"body":"bar","post_id":0}}]}}}},"posts":{"data":[{"type":"posts","id":"1","attributes":{"blog_id":0,"body":"Bar","title":"Foo"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"2","attributes":{"body":"bar","post_id":0}}]}}},{"type":"posts","id":"2","attributes":{"blog_id":0,"body":"Bas","title":"Fuubar"},"relationships":{"comments":{"data":[{"type":"comments","id":"1","attributes":{"body":"foo","post_id":0}},{"type":"comments","id":"3","attributes":{"body":"bas","post_id":0}}]}}}]}}}}
	`), blog_struct_interface)

	// data can be found in orderly structures with same signature as in original example now

	buf := bytes.NewBuffer(nil)
	jsonapi.MarshalOnePayloadEmbedded(buf, blog_struct_interface)

	fmt.Println(buf.String())
}
