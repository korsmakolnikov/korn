package main

import (
	"fmt"
	"os"

	"github.com/korsmakolnikov/kornvimgen/pkg/configuration"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(currentCmd)
}

var currentCmd = &cobra.Command{
	Use:   "current [build name]",
	Short: "set the current build name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCurrent(cmd, args)
	},
}

func executeCurrent(_ *cobra.Command, args []string) {
	targetBuildName := args[0]
	if !anyBuildName(targetBuildName, &config) {
		fmt.Println("the build you are trying to set as the current build does not exist")
		os.Exit(1)
	}

	config.CurrentBuild = args[0]
	config.Store()
	fmt.Printf("configuration set to '%s'", args[0])
}

func anyBuildName(targetBuildName string, config *configuration.Config) (res bool) {
	for k := range config.Builds {
		if k == targetBuildName {
			res = true
		}
	}
	return
}
