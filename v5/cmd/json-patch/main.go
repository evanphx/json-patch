package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	jsonpatch "github.com/evanphx/json-patch/v5"
	flags "github.com/jessevdk/go-flags"
)

type opts struct {
	PatchFilePaths []FileFlag `long:"patch-file" short:"p" value-name:"PATH" description:"Path to file with one or more operations"`
	Merge          bool       `long:"merge" short:"m" description:"Treat patches as RFC 7396 ones"`
}

func main() {
	var o opts
	_, err := flags.Parse(&o)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	patches := make([]jsonpatch.Patch, len(o.PatchFilePaths))
	patchDatas := make([][]byte, len(patches))

	for i, patchFilePath := range o.PatchFilePaths {
		var bs []byte
		bs, err = ioutil.ReadFile(patchFilePath.Path())
		if err != nil {
			log.Fatalf("error reading patch file: %s", err)
		}

		if !o.Merge {
			var patch jsonpatch.Patch
			patch, err = jsonpatch.DecodePatch(bs)
			if err != nil {
				log.Fatalf("error decoding patch file: %s", err)
			}

			patches[i] = patch
		}
		patchDatas[i] = bs
	}

	doc, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading from stdin: %s", err)
	}

	mdoc := doc
	for i, patchData := range patchDatas {
		if o.Merge {
			mdoc, err = jsonpatch.MergePatch(mdoc, patchData)
		} else {
			mdoc, err = patches[i].Apply(mdoc)
		}
		if err != nil {
			log.Fatalf("error applying patch: %s", err)
		}
	}

	fmt.Printf("%s", mdoc)
}
