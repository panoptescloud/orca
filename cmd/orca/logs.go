package main

import (
	"github.com/panoptescloud/orca/internal/controller"
	"github.com/spf13/cobra"
)

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
