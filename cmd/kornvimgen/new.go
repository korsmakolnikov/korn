package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/korsmakolnikov/kornvimgen/pkg/configuration"
	"github.com/korsmakolnikov/kornvimgen/pkg/templates"
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

func executeNew(_ *cobra.Command, args []string) {
	buildName := guessBuildName(args)

	settingPath, err := configuration.DefaultSettingPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	basePath := filepath.Join(settingPath, "builds/")
	buildPath := filepath.Join(basePath, buildName)

	if err := config.AddBuild(buildName, buildPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initTemplatable := templates.NewInit(basePath, buildName)
	if err := executeTemplatable(initTemplatable); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	customPlugin := viper.GetString("custom_plugin")
	packageLuaTemplatable := templates.NewPackage(basePath, buildName, customPlugin)
	if err := executeTemplatable(packageLuaTemplatable); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := config.Store(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeTemplatable(tmpl templates.Templetable) error {
	err := tmpl.Parse()
	if err != nil {
		return errors.Join(errors.New("[error] parsing template"), err)
	}

	if err := tmpl.Prepare(tmpl.ToPath()); err != nil {
		return errors.Join(errors.New("[error], preparing template"), err)
	}

	outputPath := tmpl.ToPath().String()
	file, err := os.OpenFile(outputPath, os.O_RDWR, 0644)
	if err != nil {
		return errors.Join(errors.New("[error] opening output file"), err)
	}
	defer file.Close()

	if err := tmpl.Execute(file); err != nil {
		return errors.Join(errors.New("[error] executing template"), err)
	}

	fmt.Println("[info] file created at path", outputPath)
	return nil
}
