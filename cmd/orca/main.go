package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

var svcContainer *services = &services{}

var rootCmd = &cobra.Command{
	Use:   "orca",
	Short: "Orca is used to orchestrate complex development environments with docker compose.",
	Long: `This can be used to manage multiple projects in different repositories
so that they can be started/stopped as a singular unit, and provides common
utilities to interact with services form anywhere on the host.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world!")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version you're using.",
	RunE:  handleVersion,
}

var sysCmd = &cobra.Command{
	Use:   "sys",
	Short: "Commands for handling the installation of this tool.",
	RunE:  handleGroup,
}

var sysCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks for dependencies.",
	Long:  `Checks that the following system dependencies are installed and usable.`,
	RunE:  handleCheck,
}

func handleGroup(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func init() {
	versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")

	sysCmd.AddCommand(sysCheckCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sysCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
