package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "start Neovim with the current \"build\"",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		executeRun(cmd, args)
	},
}

func executeRun(_ *cobra.Command, args []string) {
	var buildPath string
	for build, path := range config.Builds {
		if config.CurrentBuild == build {
			buildPath = path
		}
	}

	// Retrieving absolute Neovim executable path
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		fmt.Println("Cannot locate Neovim", err)
		os.Exit(1)
	}

	nvimArguments := buildNeovimArguments(args, buildPath)
	if err := syscall.Exec(nvimPath, nvimArguments, []string{"KORNVIM_TEST_FLAG", "1"}); err != nil {
		fmt.Println("Error running the build", err)
		os.Exit(1)
	}

	fmt.Println("Neovim execution terminated normally")
}

func buildNeovimArguments(args []string, buildPath string) []string {
	var argumentToForward []string
	argumentToForward = args
	// if len(args) > 1 {
	// 	argumentToForward = args[1:]
	// }
	configPath := filepath.Join(buildPath, "init.lua")
	runtimePath := fmt.Sprintf("set runtimepath^=%s", buildPath)
	nvimArgs := []string{"-u", configPath, "--cmd", runtimePath}
	nvimArgs = append(nvimArgs, argumentToForward...)

	return nvimArgs
}
