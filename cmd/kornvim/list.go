package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list available builds (Neovim configuration folders) in your system",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		executeList(cmd)
	},
}

func executeList(_ *cobra.Command) {
	fmt.Printf("Build\t\tPath\n")
	for k, v := range config.Builds {
		fmt.Printf("%s:\t%s\n", k, v)
	}
}
