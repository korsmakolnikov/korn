package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCMd = &cobra.Command{
	Use:   "run [build name]",
	Short: "runs the buld",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		executeRun(cmd, args)
	},
}

func executeRun(_ *cobra.Command, args []string) {
	// Retrieving absolute Neovim executable path
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		log.Fatalln("Cannot locate Neovim", err)
	}

	// Building the build path
	buildName := defaultBuildName
	if len(args) > 0 {
		buildName = args[0]
	}

	buildPath := fmt.Sprintf("./%s", buildName)

	// Building the cmd
	os.Setenv("KORNVIM_TEST_FLAG", "1")

	nvimArguments := buildNeovimArguments(args, buildPath)
	cmd := exec.Command(nvimPath, nvimArguments...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// Running the command
	err = cmd.Run()
	if err != nil {
		log.Fatalln("Error running the build", err)
	}

	log.Println("Neovim execution terminated normally")
}

func init() {
	rootCmd.AddCommand(runCMd)
}

func buildNeovimArguments(args []string, buildPath string) []string {
	var argumentToForward []string
	if len(args) > 1 {
		argumentToForward = args[1:]
	}
	configPath := fmt.Sprintf("%s/init.lua", buildPath)
	runtimePath := fmt.Sprintf("set runtimepath^=%s", buildPath)
	nvimArgs := []string{"-u", configPath, "--cmd", runtimePath}
	nvimArgs = append(nvimArgs, argumentToForward...)

	return nvimArgs
}
