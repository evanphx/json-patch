# JSON-Patch
`jsonpatch` provides the ability to decode and apply JSON patches against
documents, as well as generate merge-patches.

Specifically, this package provides the ability to apply [RFC6902 JSON patches](http://tools.ietf.org/html/rfc6902) 
as well as create [RFC7396 JSON Merge Patches](https://tools.ietf.org/html/rfc7396).

[![GoDoc](https://godoc.org/github.com/evanphx/json-patch?status.svg)](http://godoc.org/github.com/evanphx/json-patch)
[![Build Status](https://travis-ci.org/evanphx/json-patch.svg?branch=master)](https://travis-ci.org/evanphx/json-patch)
[![Report Card](https://goreportcard.com/badge/github.com/evanphx/json-patch)](https://goreportcard.com/report/github.com/evanphx/json-patch)

# Get It!

**Latest and greatest**: 
```bash
go get -u github.com/evanphx/json-patch
```

**Stable Versions**:
* Version 3: `go get -u gopkg.in/evanphx/json-patch.v3`

(previous versions below `v3` are unavailable)

# Use It!
* [Create a merge patch](#create-a-merge-patch)
* [Create and apply a Patch](#create-and-apply-a-patch)
* [Comparing JSON documents](#comparing-json-documents)

## Create a merge patch
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

## Create and apply a Patch
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

## Run It!
### CLI for comparing JSON documents
You can install the commandline program `json-patch`.

This program can take multiple JSON patch documents as arguments, 
and fed a JSON document from `stdin`. It will apply the patch(es) against 
the document and output the modified doc.

**patch.1.json**
```json
[
    {"op": "replace", "path": "/name", "value": "Jane"},
    {"op": "remove", "path": "/height"}
]
```

**patch.2.json**
```json
[
    {"op": "add", "path": "/address", "value": "123 Main St"},
    {"op": "replace", "path": "/age", "value": "21"}
]
```

**document.json**
```json
{
    "name": "John",
    "age": 24,
    "height": 3.21
}
```

You can then run:

```bash
$ go install github.com/evanphx/json-patch/cmd/json-patch
$ cat document.json | json-patch -p patch.1.json -p patch.2.json
{"address":"123 Main St","age":"21","name":"Jane"}
```

# Help It!
Contributions are welcomed! Leave [an issue](https://github.com/evanphx/json-patch/issues)
or [create a PR](https://github.com/evanphx/json-patch/compare).


Before creating a pull request, we'd ask that you make sure tests are passing
and that you have added new tests when applicable.

Contributors can run tests using:

```bash
go test -cover ./...
```

Builds for pull requests are tested automatically 
using [TravisCI](https://travis-ci.org/evanphx/json-patch).
