package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	cobra "github.com/spf13/cobra"
)

var (
	patchFiles []string
	patch      patcher = rfc6902{}
	indent     string
)

type patcher interface {
	patch(json []byte, patches [][]byte) ([]byte, error)
}

type rfc6902 struct{}

func (rfc6902) patch(json []byte, patches [][]byte) ([]byte, error) {
	var jsonPatches []jsonpatch.Patch
	for i, b := range patches {
		p, err := jsonpatch.DecodePatch(b)
		if err != nil {
			return nil, fmt.Errorf("could not decode patch #%d: %w", i, err)
		}
		jsonPatches = append(jsonPatches, p)
	}

	for i, p := range jsonPatches {
		var err error
		json, err = p.Apply(json)
		if err != nil {
			return nil, fmt.Errorf("could not apply patch #%d: %w", i, err)
		}
	}
	return json, nil
}

type rfc7396 struct{}

func (rfc7396) patch(json []byte, patches [][]byte) ([]byte, error) {
	for i, p := range patches {
		var err error
		json, err = jsonpatch.MergePatch(json, p)
		if err != nil {
			return nil, fmt.Errorf("could not apply patch #%d: %w", i, err)
		}
	}
	return json, nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "patch [flags] [JSON document]",
		Short: "Apply one or more JSON patches to a JSON document.",
		Long: `Apply one or more JSON patches to a JSON document.

The JSON document file can either be passed in as an argument, or piped through stdin.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var patches [][]byte
			for _, p := range patchFiles {
				b, err := ioutil.ReadFile(p)
				if err != nil {
					return fmt.Errorf("could not read %q: %w", p, err)
				}
				patches = append(patches, b)
			}

			var doc []byte
			var err error
			if len(args) == 0 {
				doc, err = ioutil.ReadAll(os.Stdin)
			} else {
				doc, err = ioutil.ReadFile(args[0])
			}
			if err != nil {
				return err
			}

			res, err := patch.patch(doc, patches)
			if err != nil {
				return err
			}
			if indent != "" {
				var b bytes.Buffer
				if err := json.Indent(&b, res, "", indent); err != nil {
					return err
				}
				res = b.Bytes()
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().VarP(FilesFlag{&patchFiles}, "patch-file", "p", "Path to file with one or more operations")
	cmd.Flags().VarP(patcherFlag{&patch}, "patch-format", "f", "The format of the patches. One of RFC6902 or RFC7396.")
	cmd.Flags().StringVar(&indent, "indent", "", "What indent to use when formatting the result.")

	rootCmd.AddCommand(cmd)
}

type patcherFlag struct {
	p *patcher
}

func (f patcherFlag) Set(value string) error {
	switch strings.ToUpper(value) {
	case "RFC6902":
		*f.p = rfc6902{}
		return nil
	case "RFC7396":
		*f.p = rfc7396{}
		return nil
	default:
		return fmt.Errorf("unknown patch format %q", value)
	}
}

func (f patcherFlag) String() string {
	if *f.p == nil {
		return ""
	}
	n := fmt.Sprintf("%T", *f.p)
	s := strings.Split(n, ".")
	return strings.ToUpper(s[1])
}

func (f patcherFlag) Type() string {
	return "RFC6902|RFC7396"
}
