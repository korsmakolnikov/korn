package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [build name]",
	Short: "cleanup the build directory",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		executeClean(cmd, args)
	},
}

func executeClean(_ *cobra.Command, args []string) {
	buildName := defaultBuildName
	if len(args) > 0 {
		buildName = args[0]
	}

	buildPath := fmt.Sprintf("./%s/", buildName)
	if err := os.RemoveAll(buildPath); err != nil {
		log.Fatalln("Error cleaning up the build directory")
	}
	log.Println("Cleaned up the build directory", buildPath)
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
