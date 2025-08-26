package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func handleUtilGenDocs(cmd *cobra.Command, _ []string) error {
	err := doc.GenMarkdownTree(rootCmd, "./docs/CLI")

	return err
}
