package main

import (
	"github.com/adamkirk/orca/internal/git"
	"github.com/spf13/cobra"
)

func handleGCo(cmd *cobra.Command, args []string) error {
	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	g := svcContainer.GetGit()

	searchTerm := ""

	if len(args) > 0 {
		searchTerm = args[0]
	}

	allProjects, err := cmd.Flags().GetBool("all")
	cobra.CheckErr(err)

	shouldCreate, err := cmd.Flags().GetBool("create")
	cobra.CheckErr(err)

	shouldPull, err := cmd.Flags().GetBool("pull")
	cobra.CheckErr(err)

	err = g.Checkout(git.CheckoutDTO{
		Name:        searchTerm,
		AllProjects: allProjects,
		Workspace:   ws,
		Project:     project,
		Pull:        shouldPull,
		Create:      shouldCreate,
	})

	if err != nil {
		return err
	}

	return nil
}

func handleGBranches(cmd *cobra.Command, args []string) error {
	g := svcContainer.GetGit()

	searchTerm := ""

	if len(args) > 0 {
		searchTerm = args[0]
	}

	return g.ShowBranches(git.SearchBranchesDTO{
		Search: searchTerm,
	})
}

func handleGRebaseInteractively(cmd *cobra.Command, args []string) error {
	g := svcContainer.GetGit()

	amount, err := cmd.Flags().GetInt("number")
	cobra.CheckErr(err)

	return g.RebaseInteractively(git.RebaseInteractivelyDTO{
		Amount: amount,
	})
}

func handleGPush(cmd *cobra.Command, args []string) error {
	g := svcContainer.GetGit()

	force, err := cmd.Flags().GetBool("force")
	cobra.CheckErr(err)

	return g.Push(git.PushDTO{
		Force: force,
	})
}

func handleGUndo(cmd *cobra.Command, args []string) error {
	g := svcContainer.GetGit()

	autoConfirm, err := cmd.Flags().GetBool("yes")
	cobra.CheckErr(err)

	amount, err := cmd.Flags().GetInt("number")
	cobra.CheckErr(err)

	return g.UndoLastXCommits(git.UndoLastXCommitsDTO{
		Amount:           amount,
		SkipConfirmation: autoConfirm,
	})
}

func handleGLogl(cmd *cobra.Command, args []string) error {
	g := svcContainer.GetGit()

	amount, err := cmd.Flags().GetInt("number")
	cobra.CheckErr(err)

	return g.Logl(git.LoglDTO{
		Amount: amount,
	})
}

func handleGPull(cmd *cobra.Command, args []string) error {
	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)
	allProjects, err := cmd.Flags().GetBool("all")
	cobra.CheckErr(err)

	g := svcContainer.GetGit()

	return g.Pull(git.PullDTO{
		AllProjects: allProjects,
		Workspace:   ws,
		Project:     project,
	})
}
