## JSON-Patch

Provides the abiilty to modify and test a JSON according to a
[RFC6902 JSON patch](http://tools.ietf.org/html/rfc6902).

*Version*: **1.0**


### API Usage

* Given a `[]byte`, obtain a Patch object

  `obj, err := jsonpatch.DecodePatch(patch)`

* Apply the patch and get a new document back

  `out, err := obj.Apply(doc)`

* Bonus API: compare documents for structural equality

  `jsonpatch.Equal(doca, docb)`

