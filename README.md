## JSON-Patch

Provides the ability to modify and test a JSON according to a
[RFC6902 JSON patch](http://tools.ietf.org/html/rfc6902) and [RFC7396 JSON Merge Patch](https://tools.ietf.org/html/rfc7396).

*Version*: **1.0**

[![GoDoc](https://godoc.org/github.com/evanphx/json-patch?status.svg)](http://godoc.org/github.com/evanphx/json-patch)
[![Build Status](https://travis-ci.org/evanphx/json-patch.svg?branch=master)](https://travis-ci.org/evanphx/json-patch)
[![Report Card](https://goreportcard.com/badge/github.com/evanphx/json-patch)](https://goreportcard.com/report/github.com/evanphx/json-patch)

### API Usage
#### Create a merge patch
Given both an original JSON document and a modified JSON document, you can create
a "merge patch" document, used to describe the changes needed to convert from the
original to the modified.

```go
package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func main() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	modified := []byte(`{"name": "Jane", "age": 24}`)

	patch, err := jsonpatch.CreateMergePatch(original, modified)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(patch))
}
```

When ran, you get the following output:

```bash
$ go run main.go
{"height":null,"name":"Jane"}
```

#### Create and apply a Patch
You can create patch objects using `DecodePatch([]byte)`, which can then 
be applied against JSON documents.

The following is an example of creating a patch from two operations, and
applying it against a JSON document.

```go
package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func main() {
	document := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	patchJSON := []byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`)

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(document)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(modified))
}
```

When ran, you get the following output:

```bash
$ go run main.go
{"age":24,"name":"Jane"}
```

#### Comparing JSON documents

If you have two JSON documents, you can  compare documents for structural
equality using `Equal`, which will compare the two documents based on their
actual keys and values, and not on whitespace or attribute ordering.

For example:

```go
package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func main() {
	firstDoc := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	secondDoc := []byte(`{
		"height": 3.21,
		"age"	: 24,
		"name"	: "John"
	}`)

	didMatch := jsonpatch.Equal(firstDoc, secondDoc)
	if didMatch {
		fmt.Println("The two documents were the same structurally.")
	} else {
		fmt.Println("The two documents were structurally different.")
	}
}
```

When ran, you get the following output:

```bash
$ go run main.go
The two documents were structurally equal.
```
