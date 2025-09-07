package main

import (
	"github.com/panoptescloud/orca/internal/controller"
	"github.com/spf13/cobra"
)

func handleHosts(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)

	return ctrl.Hosts(controller.HostsDTO{
		Workspace: ws,
	})
}
