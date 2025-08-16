package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/spf13/cobra"
)

func handleCheck(_ *cobra.Command, _ []string) error {
	tui := svcContainer.GetTui()

	err := hostsys.EnsureToolsExist(hostsys.EnsureToolsExistDTO{
		Tools: hostsys.RequiredTools(),
	})

	if err == nil {
		tui.Success("All required tools are present!")
		return nil
	}

	if notFoundErr, ok := err.(common.ErrToolsNotFoundOnSystem); ok {
		tui.Error(notFoundErr.Error())
		os.Exit(1)
	}

	return err
}