package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
  version string = "dev"
  commit string = "none"
  date string = "unknown"
)

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
  Use: "version",
  Short: "Show the current version you're using.",
  RunE: handleVersion,
}

var checkCmd = &cobra.Command{
  Use: "check",
  Short: "Checks for dependencies.",
  Long: `Checks that the following system dependencies are installed and usable:
- git
- docker`,
  RunE: handleCheck,
}

func init() {
  versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")

  rootCmd.AddCommand(versionCmd)
  rootCmd.AddCommand(checkCmd)
}

func main() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}