package main

import (
	"os"

	"github.com/panoptescloud/orca/internal/common"
	"github.com/panoptescloud/orca/internal/hostsys"
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

func handleSysInstall(_ *cobra.Command, args []string) error {
	tui := svcContainer.GetTui()
	if len(args) < 1 {
		return tui.RecordIfError("Must provide the name of the tool as the first argument", common.ErrArgumentRequired{
			Name:  "tool",
			Index: 1,
		})
	} else if len(args) > 1 {
		return tui.RecordIfError("Too many arguments provided!", common.ErrTooManyArguments{
			Expected: 1,
		})
	}

	sys := svcContainer.GetHostSystem()

	return sys.Install(hostsys.InstallDTO{
		Tool: hostsys.AvailableTool(args[0]),
	})
}
