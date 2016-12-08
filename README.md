# jsonstruct

[![Build Status](https://travis-ci.org/berlincount/jsonstruct.svg?branch=master)](https://travis-ci.org/berlincount/jsonstruct) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/berlincount/jsonstruct) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/berlincount/jsonstruct/master/LICENSE)


A JSON deserializer for Go structures for Go 1.7+

Also visit [Godoc](http://godoc.org/github.com/berlincount/jsonstruct).

## Installation

```
go get -u github.com/berlincount/jsonstruct
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
[sqlx](https://github.com/jmoiron/sqlx) or [GORM](https://github.com/jinzhu/gorm).

jsonstruct uses the following structures for descriptions:

```go
type Field struct {
        Name      string            "json:\"name\""
        Type      string            "json:\"type\""
        Tags      reflect.StructTag "json:\"tags\""
}

type Struct struct {
        Struct string "json:\"struct\""
        Fields []Field
}
```

which allows e.g. to describe the example structures from [JSON
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

[examples/jsonapi/jsonapi.go](https://github.com/berlincount/jsonstruct/blob/master/examples/jsonapi/jsonapi.go)

[examples/sqlx/sqlx.go](https://github.com/berlincount/jsonstruct/blob/master/examples/sqlx/sqlx.go)

[examples/gorm/gorm.go](https://github.com/berlincount/jsonstruct/blob/master/examples/gorm/gorm.go)

These runnable files show using jsonstruct with JSON API as well as in conjunction with a database using [sqlx](https://github.com/jmoiron/sqlx) or [GORM](https://github.com/jinzhu/gorm).

You can use [GB](https://getgb.io/) to build example binaries.

## Contributing

Fork, Change, Pull Request *with tests*.
