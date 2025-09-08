package commands

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

		viper.AddConfigPath(home)
		viper.SetConfigFile(".kornvimgen")
	}

	// fetch environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Configuration file loaded", viper.ConfigFileUsed())
	}
}
