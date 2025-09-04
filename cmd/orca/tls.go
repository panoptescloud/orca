package main

import (
	"github.com/panoptescloud/orca/internal/tls"
	"github.com/spf13/cobra"
)

func handleTLSGen(cmd *cobra.Command, args []string) error {
	cm := svcContainer.GetCertificateManager()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)

	return cm.Generate(tls.GenerateDTO{
		WorkspaceName: ws,
	})
}
