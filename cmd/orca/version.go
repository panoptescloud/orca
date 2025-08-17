package main

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/spf13/cobra"
)

func handleVersion(cmd *cobra.Command, _ []string, _ *config.Config) error {
	short, err := cmd.Flags().GetBool("short")

	if err != nil {
		return err
	}

	if short {
		fmt.Println(version)
		return nil
	}

	fmt.Printf("%s (%s@%s)\n", version, commit, date)

	return nil
}
