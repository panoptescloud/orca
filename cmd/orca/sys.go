package main

import (
	"os"

	"github.com/spf13/cobra"
)

func handleCheck(_ *cobra.Command, _ []string) error {
	tui := svcContainer.GetTui()

	hs := svcContainer.GetHostSystem()

	err := hs.VerifySetup()

	if err != nil {
		os.Exit(1)
		return nil
	}

	tui.Success("All requirements met!")

	return nil
}
