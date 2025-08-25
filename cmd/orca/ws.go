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

func handleWsLs(cmd *cobra.Command, args []string) error {
	manager := svcContainer.GetWorkspaceManager()

	return manager.Ls(workspaces.LsDTO{})
}

func handleWsClone(cmd *cobra.Command, args []string) error {
	manager := svcContainer.GetWorkspaceManager()

	name, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)

	projectName, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	to, err := cmd.Flags().GetString("target")
	cobra.CheckErr(err)

	return manager.Clone(workspaces.CloneDTO{
		WorkspaceName: name,
		Project:       projectName,
		To:            to,
	})
}
