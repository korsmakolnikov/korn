package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/korsmakolnikov/kornvimgen/pkg/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new [path of the custom nvim configuration folder]",
	Short: "generates a basic nvim configuration in a custom directory, installing lazy and a plugin of your choice with your custom configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeNew(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func guessNewBuildPath(buildName string) (string, error) {
	settingPath, err := configuration.DefaultSettingPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(settingPath, "builds/", buildName), nil
}

// TODO rifattorizza estraendo i template in un package
func executeNew(_ *cobra.Command, args []string) {
	buildName := guessBuildName(args)
	customPlugin := viper.GetString("custom_plugin")

	buildPath, err := guessNewBuildPath(buildName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addBuildToConfig(buildName, buildPath)

	luaFolderPath := filepath.Join(buildPath, "lua")
	fmt.Println("Creating the directory", buildPath)
	if err := os.MkdirAll(luaFolderPath, os.ModePerm); err != nil {
		fmt.Println("cannot create lua directory in the setting directory", err)
		os.Exit(1)
	}

	initFilePath := filepath.Join(luaFolderPath, "init.lua")
	initData := struct {
		BuildName string
	}{
		BuildName: buildName,
	}

	executeLuaTemplate(initFilePath, InitLuaTemplate, initData)

	packageFilePath := filepath.Join(luaFolderPath, "packages.lua")
	packageData := struct {
		CustomPlugin string
		PackageName  string
	}{
		CustomPlugin: customPlugin,
		PackageName:  guessRepoName(customPlugin),
	}
	executeLuaTemplate(packageFilePath, PackaesLuaTemplate, packageData)

	if err := config.Store(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addBuildToConfig(buildName string, buildPath string) {
	if err := config.AddBuild(buildName, buildPath); err != nil {
		os.Exit(1)
	}
}

func executeLuaTemplate(filePath string, luaTemplate string, data any) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error while creating %s file: %+v", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	tmpl, err := template.New(filePath).Parse(luaTemplate)
	if err != nil {
		fmt.Printf("Error while parsing %s file template: %+v", filePath, err)
		os.Exit(1)
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Error while execution of %s file template failed: %+v", filePath, err)
		os.Exit(1)
	}

	fmt.Printf("%s file has been successfully generated", filePath)
}

func guessRepoName(pluginNamespace string) string {
	return filepath.Base(pluginNamespace)
}
