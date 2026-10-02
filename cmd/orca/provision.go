package main

import (
	"github.com/panoptescloud/orca/internal/controller"
	"github.com/spf13/cobra"
)

func handleProvision(cmd *cobra.Command, args []string) error {
	c := svcContainer.GetController()

	return c.Provision(controller.ProvisionDTO{})
}
