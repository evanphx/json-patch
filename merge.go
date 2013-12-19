package jsonpatch

import (
	"encoding/json"
	"fmt"
)

func merge(cur, patch *lazyNode) *lazyNode {
	curDoc, err := cur.intoDoc()

	if err != nil {
		pruneNulls(patch)
		return patch
	}

	patchDoc, err := patch.intoDoc()

	if err != nil {
		return patch
	}

	mergeDocs(curDoc, patchDoc)

	return cur
}

func mergeDocs(doc, patch *partialDoc) {
	for k, v := range *patch {
		if v == nil {
			delete(*doc, k)
		} else {
			cur, ok := (*doc)[k]

			if !ok {
				pruneNulls(v)
				(*doc)[k] = v
			} else {
				(*doc)[k] = merge(cur, v)
			}
		}
	}
}

func pruneNulls(n *lazyNode) {
	sub, err := n.intoDoc()

	if err == nil {
		pruneDocNulls(sub)
	} else {
		ary, err := n.intoAry()

		if err == nil {
			pruneAryNulls(ary)
		}
	}
}

func pruneDocNulls(doc *partialDoc) *partialDoc {
	for k, v := range *doc {
		if v == nil {
			delete(*doc, k)
		} else {
			pruneNulls(v)
		}
	}

	return doc
}

func pruneAryNulls(ary *partialArray) *partialArray {
	var newAry []*lazyNode

	for _, v := range *ary {
		if v != nil {
			pruneNulls(v)
			newAry = append(newAry, v)
		}
	}

	*ary = newAry

	return ary
}

var eBadJSONDoc = fmt.Errorf("Invalid JSON Document")
var eBadJSONPatch = fmt.Errorf("Invalid JSON Patch")

func MergePatch(docData, patchData []byte) ([]byte, error) {
	doc := new(partialDoc)

	docErr := json.Unmarshal(docData, doc)

	patch := new(partialDoc)

	patchErr := json.Unmarshal(patchData, patch)

	if _, ok := docErr.(*json.SyntaxError); ok {
		return nil, eBadJSONDoc
	}

	if _, ok := patchErr.(*json.SyntaxError); ok {
		return nil, eBadJSONPatch
	}

	if docErr == nil && *doc == nil {
		return nil, eBadJSONDoc
	}

	if patchErr == nil && *patch == nil {
		return nil, eBadJSONPatch
	}

	if docErr != nil || patchErr != nil {
		// Not an error, just not a doc, so we turn straight into the patch
		if patchErr == nil {
			doc = pruneDocNulls(patch)
		} else {
			patchAry := new(partialArray)
			patchErr = json.Unmarshal(patchData, patchAry)

			if patchErr != nil {
				return nil, eBadJSONPatch
			}

			pruneAryNulls(patchAry)

			out, patchErr := json.Marshal(patchAry)

			if patchErr != nil {
				return nil, eBadJSONPatch
			}

			return out, nil
		}
	} else {
		mergeDocs(doc, patch)
	}

	return json.Marshal(doc)
}
