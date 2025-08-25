package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func handleConfigShow(cmd *cobra.Command, _ []string) error {
	contents, err := yaml.Marshal(svcContainer.GetConfig().GetRuntimeConfig())

	cobra.CheckErr(err)

	fmt.Print(string(contents))

	return nil
}

func handleConfigPath(cmd *cobra.Command, _ []string) error {
	fmt.Println(getConfigFilePath())

	return nil
}
