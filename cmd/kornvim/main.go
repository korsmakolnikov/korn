package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

type launchingArguments struct {
	executable string
	buildName  string
}

func (l *launchingArguments) args() ([]string, error) {
	return []string{
		"nvim",
		"-u",
		fmt.Sprintf("%s/init.lua", l.buildName),
		"--cmd",
		fmt.Sprintf("set runtimepath^=%s", l.buildName),
	}, nil
}

func (l *launchingArguments) exePath() string {
	return l.executable
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	configDir := fmt.Sprintf("%s/.config/kornvim", home)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	_ = fmt.Sprintf("%s/config", configDir)

	// TODO creare il file se non esiste
	// TODO leggere e scrivere la configurazione

	cmd := launchingArguments{
		executable: "/usr/bin/nvim",
		buildName:  "kornvim",
	}

	args, err := cmd.args()
	if err != nil {
		log.Fatalln(err)
	}

	args = append(args, os.Args[1:]...)

	if err := syscall.Exec(cmd.exePath(), args, os.Environ()); err != nil {
		log.Fatalln(err)
	}
}
