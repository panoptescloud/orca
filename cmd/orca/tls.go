package main

import (
	"github.com/adamkirk/orca/internal/tls"
	"github.com/spf13/cobra"
)

var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "Commands to do with TLS certificates.",
	RunE:  handleGroup,
}

var tlsGenCmd = &cobra.Command{
	Use: "gen",
	Short: `Generates TLS certificates for the current workspace.
If no root certificate has been generated, one will be generated.`,
	Run: errorHandlerWrapper(handleTLSGen, 1),
}

func init() {
	addWorkspaceOption(tlsGenCmd, false)
	tlsCmd.AddCommand(tlsGenCmd)
	rootCmd.AddCommand(tlsCmd)
}

func handleTLSGen(cmd *cobra.Command, args []string) error {
	cm := svcContainer.GetCertificateManager()

	ws, err := cmd.Flags().GetString("workspace")
	cobra.CheckErr(err)

	return cm.Generate(tls.GenerateDTO{
		WorkspaceName: ws,
	})
}
