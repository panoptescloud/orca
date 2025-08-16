package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func handleCheck(_ *cobra.Command, _ []string) error {
	fmt.Println("Checking docker is executable...")

	return nil
}