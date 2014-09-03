## JSON-Patch

Provides the abiilty to modify and test a JSON according to a
[RFC6902 JSON patch](http://tools.ietf.org/html/rfc6902) and [JSON Merge Patch](http://tools.ietf.org/html/draft-ietf-appsawg-json-merge-patch-07).

*Version*: **1.0**


### API Usage

* Given a `[]byte`, obtain a Patch object

  `obj, err := jsonpatch.DecodePatch(patch)`

* Apply the patch and get a new document back

  `out, err := obj.Apply(doc)`

* Create a JSON Merge Patch document based on two json documents (a to b):

  `mergeDoc, err := jsonpatch.CreateMergePatch(a, b)`
 
* Bonus API: compare documents for structural equality

  `jsonpatch.Equal(doca, docb)`

