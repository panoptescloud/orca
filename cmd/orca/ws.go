package main

import (
	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/workspaces"
	"github.com/spf13/cobra"
)

func handleWsSwitch(cmd *cobra.Command, args []string) error {
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

func handleWsInit(cmd *cobra.Command, args []string) error {
	manager := svcContainer.GetWorkspaceManager()

	source, err := cmd.Flags().GetString("source")
	cobra.CheckErr(err)
	target, err := cmd.Flags().GetString("target")
	cobra.CheckErr(err)
	configFile, err := cmd.Flags().GetString("config")
	cobra.CheckErr(err)

	return manager.Initialise(workspaces.InitialiseDTO{
		SourceDirectory:   source,
		WorkspaceFileName: configFile,
		Into:              target,
	})
}
