package configuration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CurrentBuild string            `mapstructure:"current"`
	Builds       map[string]string `mapstructure:"builds"`
}

const (
	DEFAULT_CONFIG_FOLDER    = ".config/kornvim/"
	DEFAULT_CONFIG_FILE_NAME = "config"
)

func Default() Config {
	builds := make(map[string]string)

	return Config{
		CurrentBuild: "",
		Builds:       builds,
	}
}

func DefaultSettingPath() (defaultSettingPath string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	defaultSettingPath = filepath.Join(home, DEFAULT_CONFIG_FOLDER)
	return defaultSettingPath, nil
}

func DefaultSettingFilePath() (string, error) {
	defaultSettingPath, err := DefaultSettingPath()
	if err != nil {
		return "", err
	}
	defaultSettingFilePath := filepath.Join(defaultSettingPath, DEFAULT_CONFIG_FILE_NAME)
	return defaultSettingFilePath, nil
}

func UpsertConfigurationFile() error {
	settingFilePath := viper.GetString("setting")
	viper.SetConfigName("config")
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(settingFilePath)
	viper.AddConfigPath(filepath.Dir(settingFilePath))
	viper.AutomaticEnv()

	if _, err := os.Stat(settingFilePath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(settingFilePath), os.ModePerm)
		fmt.Println("[info] creating new config file at: ", settingFilePath)
		f, err := os.OpenFile(settingFilePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return errors.Join(errors.New("failing creating the configuration file"), err)
		}
		defer f.Close()

		defaultCfg := Default()
		defaultCfg.Store()
	}

	return nil
}

func (cfg Config) Store() error {
	settingFilePath := viper.ConfigFileUsed()
	viper.Set("current", cfg.CurrentBuild)
	viper.Set("builds", cfg.Builds)
	if err := viper.WriteConfigAs(settingFilePath); err != nil {
		return errors.Join(errors.New("error while storing the kornvim configuration file"), err)
	}

	return nil
}

func (cfg *Config) Load() error {
	settingFilePath := viper.GetString("setting")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.Join(fmt.Errorf("the configuration file you have provided doesn't exist at: %s", settingFilePath), err)
		} else {
			return errors.Join(fmt.Errorf("the configuration file is corrupted at: %s", settingFilePath), err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return errors.Join(errors.New("the configuration cannot be unmarshaled"), err)
	}

	if cfg.Builds == nil {
		cfg.Builds = make(map[string]string)
	}

	return nil
}

func (cfg *Config) AddBuild(buildName string, buildPath string) error {
	if _, ok := cfg.Builds[buildName]; ok {
		return errors.New("the build already exists")
	}
	cfg.Builds[buildName] = buildPath

	return nil
}

func (cfg *Config) DeleteBuild(buildName string) error {
	if cfg.CurrentBuild == buildName {
		return errors.New("the current build cannot be deleted")
	}

	delete(cfg.Builds, buildName)

	return nil
}

func (cfg Config) GetBuildPath(buildName string) (path string, err error) {
	err = fmt.Errorf("cannot find the build name '%s' in %v", buildName, cfg.GetBuilds())

	for buildNameKey, pathValue := range cfg.Builds {
		if buildName == buildNameKey {
			path = pathValue
			err = nil
			break
		}
	}

	return
}

func (cfg Config) GetCurrentPath() (string, error) {
	path, err := cfg.GetBuildPath(cfg.CurrentBuild)
	if err != nil {
		return "", errors.Join(errors.New("current build not set:"), err)
	}

	return path, nil
}

func (cfg Config) GetBuilds() (builds []string) {
	for k := range cfg.Builds {
		builds = append(builds, k)
	}

	return
}
