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
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		fmt.Printf("[error] cannot locate Neovim %s", err)
		os.Exit(1)
	}
	if len(args) == 0 {
		currentPath, err := config.GetCurrentPath()
		if err != nil {
			fmt.Printf("[error] %s", err)
			os.Exit(1)
		}
		nvimArgs := buildNeovimArguments(currentPath)
		run(nvimArgs, nvimPath)
		return
	}

	var currentBuildNotFound, buildNotFound bool
	buildPath, errBuildPath := config.GetBuildPath(args[0])
	if errBuildPath != nil {
		buildNotFound = true
	}

	currentPath, errCurrentBuildPath := config.GetCurrentPath()
	if errCurrentBuildPath != nil {
		currentBuildNotFound = true
	}
	if buildNotFound && currentBuildNotFound {
		fmt.Printf("[error] the build doesn't exists: %s, %s", errBuildPath, errCurrentBuildPath)
		os.Exit(1)
	}

	path := currentPath
	if buildPath != "" {
		path = buildPath
	}

	extraArgs := argsToForward(args, !buildNotFound)
	nvimArgs := buildNeovimArguments(path)
	nvimArgs = append(nvimArgs, extraArgs...)

	run(nvimArgs, nvimPath)
}

func run(args []string, nvimPath string) {
	env := os.Environ()
	if err := syscall.Exec(nvimPath, args, env); err != nil {
		fmt.Printf("[error] running the build %s", err)
		os.Exit(1)
	}

	fmt.Println("[info] Neovim execution terminated normally")
}

func buildNeovimArguments(buildPath string) []string {
	nvimArgs := newNvimArgs(
		filepath.Join(buildPath, "lua/init.lua"),
		fmt.Sprintf("set runtimepath^=%s", buildPath),
	).toSyscallArgs()

	return nvimArgs
}

func argsToForward(args []string, buildFound bool) []string {
	if !buildFound {
		return args
	}

	return args[1:]
}

type nvimArgs struct {
	Entrypoint  string
	Runtimepath string
}

func newNvimArgs(entrypoint string, runtimepath string) nvimArgs {
	return nvimArgs{
		Entrypoint:  entrypoint,
		Runtimepath: runtimepath,
	}
}

func (nvmargs nvimArgs) toSyscallArgs() []string {
	return []string{
		"nvim",
		"-u",
		nvmargs.Entrypoint,
		"--cmd",
		nvmargs.Runtimepath,
	}
}
