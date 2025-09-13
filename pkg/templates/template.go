package templates

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

//go:embed init.lua.tmpl
var initLuaTemplate string

//go:embed packages.lua.tmpl
var packageLuaTemplate string

var templates = map[string]string{
	"initLuaTemplate":    initLuaTemplate,
	"packageLuaTemplate": packageLuaTemplate,
}

// TODO rimuovi error
func NewInit(basePath string, buildName string) Templetable {
	name := "init.lua"
	buildPath := filepath.Join(basePath, buildName)
	outputPath := filepath.Join(basePath, buildName, "lua", name)
	return &InitTemplate{
		template:   "initLuaTemplate",
		name:       name,
		outPutPath: outputPath,
		buildName:  buildName,
		buildPath:  buildPath,
	}
}

func NewPackage(basePath string, buildName string, configurationPlugin string) Templetable {
	name := "package.lua"
	buildPath := filepath.Join(basePath, buildName)
	outputPath := filepath.Join(basePath, buildName, "lua", name)
	return &PackageTemplate{
		template:            "packageLuaTemplate",
		name:                name,
		outPutPath:          outputPath,
		buildPath:           buildPath,
		configurationPlugin: configurationPlugin,
	}
}

type Templetable interface {
	Parse() error
	Execute(file io.Writer) error
	Prepare(UpsertableIO) error
	ToPath() UpsertableIO
}

type InitTemplate struct {
	template   string
	name       string
	outPutPath string
	buildPath  string
	buildName  string
	body       *template.Template
}

func (t *InitTemplate) Parse() error {
	if tmpl, ok := templates[t.template]; ok {
		tmpl, err := template.New(t.name).Parse(tmpl)
		if err != nil {
			return errors.Join(errors.New("Error creating a new template"), err)
		}

		t.body = tmpl
		return nil
	}

	return errors.New("template not present in the database")
}

func (t InitTemplate) ToPath() UpsertableIO {
	return &UpsertableFile{path: t.outPutPath}
}

func (t *InitTemplate) Prepare(templatePath UpsertableIO) error {
	err := templatePath.touch()
	if err != nil {
		return err
	}

	return nil
}

func (t *InitTemplate) Execute(file io.Writer) error {
	args := struct {
		BuildName string
		BuildPath string
	}{
		BuildName: t.buildName,
		BuildPath: t.buildPath,
	}
	if err := t.body.Execute(file, args); err != nil {
		return fmt.Errorf("Error while execution of %s file template failed: %+v", t.outPutPath, err)
	}
	return nil
}

type PackageTemplate struct {
	template            string
	name                string
	outPutPath          string
	buildPath           string
	buildName           string
	configurationPlugin string
	body                *template.Template
}

func (t *PackageTemplate) Parse() error {
	if tmpl, ok := templates[t.template]; ok {
		tmpl, err := template.New(t.name).Parse(tmpl)
		if err != nil {
			return errors.Join(errors.New("Error creating a new template"), err)
		}

		t.body = tmpl
		return nil
	}

	return errors.New("template not present in the database")
}

func (t PackageTemplate) ToPath() UpsertableIO {
	return &UpsertableFile{path: t.outPutPath}
}

func (t *PackageTemplate) Prepare(templatePath UpsertableIO) error {
	err := templatePath.touch()
	if err != nil {
		return err
	}

	return nil
}

func (t *PackageTemplate) Execute(file io.Writer) error {
	args := struct {
		BuildName           string
		BuildPath           string
		ConfigurationPlugin string
		PackageName         string
	}{
		BuildName:           t.buildName,
		BuildPath:           t.buildPath,
		ConfigurationPlugin: t.configurationPlugin,
		PackageName:         guessRepoName(t.configurationPlugin),
	}
	if err := t.body.Execute(file, args); err != nil {
		return fmt.Errorf("Error while execution of %s file template failed: %+v", t.outPutPath, err)
	}
	return nil
}

func guessRepoName(pluginNamespace string) string {
	return filepath.Base(pluginNamespace)
}

type UpsertableIO interface {
	touch() error
	fmt.Stringer
}

type UpsertableFile struct {
	path string
}

func (u *UpsertableFile) String() string {
	return u.path
}

func (u *UpsertableFile) touch() error {
	folder := filepath.Dir(u.path)
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}

	handle, err := os.Create(u.path)
	if err != nil {
		return err
	}
	defer handle.Close()

	return nil
}
