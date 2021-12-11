package main

import (
	"fmt"
	"io/ioutil"

	jsonpatch "github.com/evanphx/json-patch"
	cobra "github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:                   "create <old JSON file> <new JSON file>",
		DisableFlagsInUseLine: true,
		Short:                 "Create an RFC7396 merge patch from the difference between two JSON documents.",
		Args:                  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			orig, err := ioutil.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("could not read file %q: %w", args[0], err)
			}

			modi, err := ioutil.ReadFile(args[1])
			if err != nil {
				return fmt.Errorf("could not read file %q: %w", args[1], err)
			}

			p, err := jsonpatch.CreateMergePatch(orig, modi)
			if err != nil {
				return err
			}

			fmt.Println(string(p))
			return nil
		},
	})
}
