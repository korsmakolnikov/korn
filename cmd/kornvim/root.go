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
korvim handle a Neovim configuration database for you. You can set one "build" 
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
	viper.BindPFlag("setting_file", rootCmd.PersistentFlags().Lookup("setting"))
}

func initConfig() {
	configuration.UpsertConfigurationFile(settingFilePath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("the configuration file you have provided doesn't exist")
			os.Exit(1)
		} else {
			fmt.Println("path: ", settingFilePath)
			fmt.Println("the configuration file is corrupted", err)
			os.Exit(1)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("the configuration cannot be unmarshaled", err)
		os.Exit(1)
	}
}
