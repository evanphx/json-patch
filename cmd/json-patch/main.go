package main

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "json-patch",
	Short: "A tool for RFC6902 and RFC7396 JSON manipulation.",
}

func main() {
	// TODO: Hide the completion command once we're able. This capability was only added recently: https://github.com/spf13/cobra/issues/1507
	// rootCmd.CompletionOptions.HiddenDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
