package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: `Tails logs from docker compose project.`,
	Long:  `Basically a 'docker compose logs -f', you may supply a service to tail specifically.`,
	Run:   errorHandlerWrapper(handleLogs, 1),
}

func init() {
	addWorkspaceOption(logsCmd, false)
	addProjectOption(logsCmd)
	addServiceOption(logsCmd, false)
	rootCmd.AddCommand(logsCmd)
}

func handleLogs(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)
	service, err := cmd.Flags().GetString("service")
	cobra.CheckErr(err)

	return ctrl.Logs(controller.LogsDTO{
		Workspace: ws,
		Project:   project,
		Service:   service,
	})
}
