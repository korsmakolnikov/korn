package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [build name]",
	Short: "cleanup the build directory",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeClean(cmd, args)
	},
}

func executeClean(_ *cobra.Command, args []string) {
	buildName := guessBuildName(args)
	buildPath, err := config.GetBuildPath(buildName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.RemoveAll(buildPath); err != nil {
		fmt.Println("Error cleaning up the build directory:", err)
		os.Exit(1)
	}
	cleanBuildFromConfig(buildName)
	fmt.Println("Cleaned up the build directory", buildPath)
}

func cleanBuildFromConfig(buildName string) {
	err := config.DeleteBuild(buildName)
	if err != nil {
		fmt.Println("cannot update the configuration: ", err)
		os.Exit(1)
	}
	if err := config.Store(); err != nil {
		fmt.Println("error while attempting to update the configuration file", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
