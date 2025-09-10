package configuration

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CurrentBuild string            `mapstructure:"current_build"`
	Builds       map[string]string `mapstructure:"builds"`
}

func Default() Config {
	builds := make(map[string]string)
	builds["kornvim"] = "./kornvim"

	return Config{
		CurrentBuild: "kornvim",
		Builds:       builds,
	}
}

func DefaultSettingFilePath() (string, error) {
	var defaultSettingFilePath string

	home, err := os.UserHomeDir()
	if err != nil {
		return defaultSettingFilePath, err
	}
	defaultSettingFilePath = filepath.Join(home, ".config/kornvim/config")
	return defaultSettingFilePath, nil
}

func UpsertConfigurationFile(settingFilePath string) error {
	viper.SetConfigName("config")
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(settingFilePath)
	viper.AddConfigPath(filepath.Dir(settingFilePath))

	if _, err := os.Stat(settingFilePath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(settingFilePath), os.ModePerm)
		log.Println("creating new config file at: ", settingFilePath)
		f, err := os.OpenFile(settingFilePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Println("failing creating the configuration file due to: ", err)
			return err
		}
		defer f.Close()

		defaultCfg := Default()
		defaultCfg.Store(viper.ConfigFileUsed())
	}

	return nil
}

func (cfg Config) Store(settingFilePath string) error {
	viper.Set("current", cfg.CurrentBuild)
	viper.Set("builds", cfg.Builds)
	if err := viper.WriteConfigAs(settingFilePath); err != nil {
		log.Println("error while storing the kornvim configuration file", err)
		return err
	}

	return nil
}
