package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Starts the workspace or project.",
	Long:  `...TBD...`,
	Run:   errorHandlerWrapper(handleUp, 1),
}

func init() {
	addWorkspaceOption(upCmd, false)
	addProjectOption(upCmd)

	rootCmd.AddCommand(upCmd)
}

func handleUp(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	return ctrl.Up(controller.UpDTO{
		Workspace: ws,
		Project:   project,
	})
}
