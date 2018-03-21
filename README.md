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

#### Decode a Patch


* Given a `[]byte`, obtain a Patch object

  `obj, err := jsonpatch.DecodePatch(patch)`

* Apply the patch and get a new document back

  `out, err := obj.Apply(doc)`

* Create a JSON Merge Patch document based on two json documents (a to b):

  `mergeDoc, err := jsonpatch.CreateMergePatch(a, b)`
 
* Bonus API: compare documents for structural equality

  `jsonpatch.Equal(doca, docb)`

