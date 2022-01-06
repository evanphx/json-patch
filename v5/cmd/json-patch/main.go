package main

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "json-patch",
	Short: "A tool for RFC6902 and RFC7396 JSON manipulation.",
	Long: `A tool for RFC6902 and RFC7396 JSON manipulation.

For reverse compatiblity with a previous version of this tool, invoking this tool without any subcommand is the same as` + " `apply -f RFC6902`.",
	CompletionOptions: cobra.CompletionOptions{
		// TODO: https://github.com/spf13/cobra/issues/1507
		// Set HiddenDefaultCmd = true and clear DisableDefaultCmd.
		DisableDefaultCmd: true,
	},
}

func init() {
	var patchFiles []string

	// Add RunE to the root command for reverse compatibility. We should probably deprecate this and encourage users to use apply directly.
	rootCmd.RunE =
		func(cmd *cobra.Command, args []string) error {
			return apply(cmd, args, patchFiles, rfc6902{}, "")
		}
	rootCmd.Flags().VarP(FilesFlag{&patchFiles}, "patch-file", "p", "Path to file with one or more operations")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
