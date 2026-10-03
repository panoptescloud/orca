package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Runs a command inside one of the containers in the environment",
	Long:  `...TBD...`,
	Run:   errorHandlerWrapper(handleExec, 1),
}

func init() {
	addWorkspaceOption(execCmd, false)
	addProjectOption(execCmd)
	addServiceOption(execCmd, true)
	rootCmd.AddCommand(execCmd)
}

func handleExec(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)
	service, err := cmd.Flags().GetString("service")
	cobra.CheckErr(err)

	return ctrl.ExecOrRun(controller.ExecDTO{
		Workspace: ws,
		Project:   project,
		Service:   service,
		Args:      args,
	})
}
