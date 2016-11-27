# jsonstruct

[![Build Status](https://travis-ci.org/berlincount/go-jsonstruct.svg?branch=master)](https://travis-ci.org/berlincount/go-jsonstruct)

A JSON deserializer for Go structures for Go 1.7+

Also visit [Godoc](http://godoc.org/github.com/berlincount/go-jsonstruct).

## Installation

```
go get -u github.com/berlincount/go-jsonstruct
```

## Background

You are using [Google's JSONAPI](http://godoc.org/github.com/google/jsonapi)
package with your Go web application and have a lot of structs in your database
schema that you don't want to also have to implement as individual Go
structures.

## Introduction

jsonstruct uses [StructOf](http://golang.org/pkg/reflect/#StructOf) to
construct a [Type](http://golang.org/pkg/reflect/#Type) which can be used to
create [Value](http://golang.org/pkg/reflect/#Value)s which then can be used by
other packages using reflection for structure discovery, like
[SQLx](https://github.com/jmoiron/sqlx)

jsonstruct uses the following structures for descriptions:

```go
type Field struct {
        Name string            "json:name"
        Type string            "json:type"
        Tags reflect.StructTag "json:tags"
}
type Struct struct {
        Struct string "json:struct"
        Fields []Field
}
```

which allows to describe the example structures from [JSON
API](http://godoc.org/github.com/google/jsonapi) using the following structure:

```javascript
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
```

## Example Apps

[src/jsonapi/example.go](https://github.com/berlincount/go-jsonstruct/blob/master/src/jsonapi/example.go)

[src/database/example.go](https://github.com/berlincount/go-jsonstruct/blob/master/src/database/example.go)

These runnable files show using jsonstruct with JSON API as well as in conjunction with a database using SQLx.

You can use [GB](https://getgb.io/) to build example binaries.

## Contributing

Fork, Change, Pull Request *with tests*.
