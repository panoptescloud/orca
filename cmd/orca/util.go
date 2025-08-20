package main

import (
	"github.com/panoptescloud/orca/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func handleUtilGenDocs(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	err := doc.GenMarkdownTree(rootCmd, "./docs/CLI")

	return err
}
