package main

import (
	"fmt"
	"log"
	"os"

	"github.com/korsmakolnikov/kornvimgen/pkg/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settingFilePath string

var config configuration.Config

var rootCmd = &cobra.Command{
	Use:   "kornvim",
	Short: "A handy cli tool for running Neovim with custom configurations",
	Long: `
kornvim handle a Neovim configuration database for you. You can set one "build" 
as the current configuration and kornvim will run it as the default, or you can 
try or run different instances of Neovim pointing to different configurations.
`,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func init() {
	defaultSettingFilePath, err := configuration.DefaultSettingFilePath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().
		StringVar(&settingFilePath, "setting", defaultSettingFilePath, fmt.Sprintf("setting file (default is %s)", defaultSettingFilePath))
	viper.BindPFlag("setting", rootCmd.PersistentFlags().Lookup("setting"))
}

func initConfig() {
	configuration.UpsertConfigurationFile()
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}
}
