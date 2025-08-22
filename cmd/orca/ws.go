package main

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/config"
	"github.com/spf13/cobra"
)

func handleWsSwitch(cmd *cobra.Command, args []string, cfg *config.Config) error {
	tui := svcContainer.GetTui()

	if len(args) < 1 {
		tui.Error("Must supply a workspace!")

		return common.ErrUnknownWorkspace{
			Workspace: "",
		}
	}

	manager := svcContainer.GetConfigManager()

	return manager.SwitchWorkspace(args[0])
}
