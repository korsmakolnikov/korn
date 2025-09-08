package main

import (
	"log"

	_ "embed"

	cmd "github.com/korsmakolnikov/kornvimgen/cmd/kornvimgen/commands"
)

//go:embed templates/init.lua.tmpl
var InitLuaTemplate string

//go:embed templates/packages.lua.tmpl
var PackaesLuaTemplate string

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
