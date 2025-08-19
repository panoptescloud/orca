package main

import (
	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/git"
	"github.com/spf13/cobra"
)

func handleGCo(cmd *cobra.Command, args []string, cfg *config.Config) error {
	g := svcContainer.GetGit()

	searchTerm := ""

	if len(args) > 0 {
		searchTerm = args[0]
	}

	err := g.Checkout(git.CheckoutDTO{
		Name: searchTerm,
	})

	if err != nil {
		return err
	}

	shouldPull, err := cmd.Flags().GetBool("pull")
	cobra.CheckErr(err)

	if shouldPull {
		return g.PullCurrentBranch()
	}

	return nil
}

func handleGBranches(cmd *cobra.Command, args []string, cfg *config.Config) error {
	g := svcContainer.GetGit()

	searchTerm := ""

	if len(args) > 0 {
		searchTerm = args[0]
	}

	return g.ShowBranches(git.SearchBranchesDTO{
		Search: searchTerm,
	})
}

func handleGPull(cmd *cobra.Command, args []string, cfg *config.Config) error {
	g := svcContainer.GetGit()

	return g.PullCurrentBranch()
}
