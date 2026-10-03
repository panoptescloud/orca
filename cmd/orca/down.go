package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops the workspace or project.",
	Long:  `...TBD...`,
	Run:   errorHandlerWrapper(handleDown, 1),
}

func init() {
	addWorkspaceOption(downCmd, false)
	addProjectOption(downCmd)

	rootCmd.AddCommand(downCmd)
}

func handleDown(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	return ctrl.Down(controller.DownDTO{
		Workspace: ws,
		Project:   project,
	})
}
