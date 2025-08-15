package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

func main() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}