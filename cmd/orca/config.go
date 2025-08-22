package main

import (
	"fmt"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func handleConfigShow(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	contents, err := yaml.Marshal(cfg)

	cobra.CheckErr(err)

	fmt.Print(string(contents))

	return nil
}

func handleConfigPath(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	fmt.Println(getConfigFilePath())

	return nil
}
