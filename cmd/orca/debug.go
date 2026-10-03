package main

import (
	"github.com/adamkirk/orca/internal/controller"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Commands to aid in debugging or understanding whats happening.",
	RunE:  handleGroup,
}

var debugShowComposeConfigCmd = &cobra.Command{
	Use:   "show-compose-config",
	Short: `Shows the full generated compose config that will be used for this project.`,
	Long: `Similar to the 'docker compose config' command (it uses this under the 
hood), it will show you the resulting config after all files are merged together, 
including any relevant values from environment variables, or profile changes etc.`,
	Run: errorHandlerWrapper(handleDebugShowComposeConfig, 1),
}

var debugShowComposeCommandCmd = &cobra.Command{
	Use:   "show-compose-command",
	Short: `Shows the compose command being used for this project.`,
	Long:  `Can be used to run generic commands via: ` + "`$(orca debug show-compose-command) ps`",
	Run:   errorHandlerWrapper(handleDebugShowComposeCommand, 1),
}

func init() {
	addWorkspaceOption(debugShowComposeConfigCmd, false)
	addProjectOption(debugShowComposeConfigCmd)
	debugCmd.AddCommand(debugShowComposeConfigCmd)

	addWorkspaceOption(debugShowComposeCommandCmd, false)
	addProjectOption(debugShowComposeCommandCmd)
	debugCmd.AddCommand(debugShowComposeCommandCmd)
	rootCmd.AddCommand(debugCmd)
}

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

func handleDebugShowComposeCommand(cmd *cobra.Command, args []string) error {
	ctrl := svcContainer.GetController()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)
	project, err := cmd.Flags().GetString("project")
	cobra.CheckErr(err)

	return ctrl.ShowComposeCommand(controller.ShowComposeCommandDTO{
		Workspace: ws,
		Project:   project,
	})
}
