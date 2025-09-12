package main

import (
	_ "embed"
	"log"
)

//go:embed templates/init.lua.tmpl
var InitLuaTemplate string

//go:embed templates/packages.lua.tmpl
var PackaesLuaTemplate string

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
