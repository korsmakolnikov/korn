package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	settingFile               string
	customConfigurationPlugin string
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&settingFile, "setting", "$HOME/.kornvimgen", "setting file (default is $HOME/.kornvimgen)")
	rootCmd.PersistentFlags().StringVar(&customConfigurationPlugin, "plugin", defaultCustomPlugin, "your configuration code lazy-compatible plugin namespace")
	viper.BindPFlag("setting_file", rootCmd.PersistentFlags().Lookup("setting"))
	viper.BindPFlag("custom_plugin", rootCmd.PersistentFlags().Lookup("plugin"))
}

func initConfig() {
	// assess the setting file or pick the default one
	if settingFile != "" {
		viper.SetConfigFile(settingFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigFile(".kornvimgen")
		viper.SetDefault("current", "kornvim_test")
		viper.SetDefault("builds", map[string]string{"default": "./kornvim_test"})
	}

	// fetch environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("configuration file not found. I will create one for you")

			if err := viper.SafeWriteConfigAs(viper.ConfigFileUsed()); err != nil {
				log.Fatalln("error creating default configuration file")
			}
		} else {
			log.Fatalln("configuration file error", err)
		}

		log.Println("Configuration file loaded", viper.ConfigFileUsed())
	}
}
