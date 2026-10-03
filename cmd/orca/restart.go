package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Alias for running down && up.",
	Long:  `...TBD...`,
	Run:   errorHandlerWrapper(handleRestart, 1),
}

func init() {
	addWorkspaceOption(restartCmd, false)
	addProjectOption(restartCmd)

	rootCmd.AddCommand(restartCmd)
}

func handleRestart(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	err = ctrl.Down(controller.DownDTO{
		Workspace: ws,
		Project:   project,
	})

	if err != nil {
		return err
	}

	return ctrl.Up(controller.UpDTO{
		Workspace: ws,
		Project:   project,
	})
}
