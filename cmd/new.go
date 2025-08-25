package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultBuildName = "kornvim"
)

var newCmd = &cobra.Command{
	Use:   "new [path of the custom nvim configuration folder]",
	Short: "generates a basic nvim configuration in a custom directory, installing lazy and a plugin of your choice with your custom configuration",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		executeNew(cmd, args)
	},
}

func executeNew(_ *cobra.Command, args []string) {
	buildName := defaultBuildName
	if len(args) > 0 {
		buildName = args[0]
	}

	customPlugin := viper.GetString("custom_plugin")

	// creates the build directory and the lua child directory
	buildPath := buildPath(buildName)
	luaFolderPath := fmt.Sprintf("%s/lua/", buildPath)
	log.Println("Creating the directory", buildPath)
	if err := os.MkdirAll(luaFolderPath, os.ModePerm); err != nil {
		log.Fatalln("cannot create lua directory in the setting directory", err)
	}

	initFilePath := fmt.Sprintf("%s/init.lua", buildPath)
	initData := struct {
		BuildName string
	}{
		BuildName: buildName,
	}

	templatesPath := "./template"
	executeLuaTemplate(initFilePath, fmt.Sprintf("%s/init.lua.tmpl", templatesPath), initData)

	packageFilePath := fmt.Sprintf("%s/packages.lua", luaFolderPath)
	packageData := struct {
		CustomPlugin string
		PackageName  string
	}{
		CustomPlugin: customPlugin,
		PackageName:  guessRepoName(customPlugin),
	}
	executeLuaTemplate(packageFilePath, fmt.Sprintf("%s/packages.lua.tmpl", templatesPath), packageData)
}

func buildPath(buildName string) string {
	return fmt.Sprintf("./%s/", buildName)
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func executeLuaTemplate(filePath string, luaTemplatePath string, data any) {
	luaTemplate, err := os.ReadFile(luaTemplatePath)
	if err != nil {
		log.Fatalf("Error while reading %s template file: %+v", luaTemplatePath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error while creating %s file: %+v", filePath, err)
	}
	defer file.Close()

	tmpl, err := template.New(filePath).Parse(string(luaTemplate))
	if err != nil {
		log.Fatalf("Error while parsing %s file template: %+v", filePath, err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("Error while execution of %s file template failed: %+v", filePath, err)
	}

	log.Printf("%s file has been successfully generated", filePath)
}

func guessRepoName(pluginNamespace string) string {
	return filepath.Base(pluginNamespace)
}
