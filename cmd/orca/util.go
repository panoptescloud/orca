package main

import (
	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/updater"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func handleUtilGenDocs(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	err := doc.GenMarkdownTree(rootCmd, "./docs/CLI")

	return err
}

func handleUtilSelfUpdate(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	selfUpdater := svcContainer.GetSelfUpdater()

	to, err := cmd.Flags().GetString("to")
	cobra.CheckErr(err)

	strategy := updater.VersioningStrategyLatest

	if to != "" {
		strategy = updater.VersioningStrategySpecific
	}
	return selfUpdater.Apply(updater.ApplyDTO{
		Strategy:         strategy,
		SpecifiedVersion: to,
	})
}
