package main

import (
	"github.com/panoptescloud/orca/internal/controller"
	"github.com/spf13/cobra"
)

func handleDebugShowComposeConfig(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	return ctrl.ShowComposeConfig(controller.ShowComposeConfigDTO{
		Workspace: ws,
		Project:   project,
	})
}
