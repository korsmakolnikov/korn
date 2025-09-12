package main

import (
	"fmt"
	"log"
	"os"

	"github.com/korsmakolnikov/kornvimgen/pkg/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	settingFilePath           string
	customConfigurationPlugin string
	config                    configuration.Config
)

const defaultCustomPlugin = "korsmakolnikov/kornvim_configurator"

var rootCmd = &cobra.Command{
	Use:   "kornvimgen",
	Short: "A cli tool for generating neovim configurations",
	Long: `
kornvimgen is a tool to generate neovim configurations in a kornvim fashion
The new command generate a Neovim configuration folder (the build).`,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
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
		StringVar(
			&settingFilePath,
			"setting",
			defaultSettingFilePath,
			fmt.Sprintf("setting file (default is %s)", defaultSettingFilePath),
		)
	rootCmd.PersistentFlags().
		StringVar(
			&customConfigurationPlugin,
			"plugin",
			defaultCustomPlugin,
			"your configuration code lazy-compatible plugin namespace",
		)
	viper.BindPFlag("setting", rootCmd.PersistentFlags().Lookup("setting"))
	viper.BindPFlag("custom_plugin", rootCmd.PersistentFlags().Lookup("plugin"))
}

func initConfig() {
	configuration.UpsertConfigurationFile()
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}
}

func guessBuildName(args []string) (buildName string) {
	if len(args) > 0 {
		buildName = args[0]
	}

	return
}
