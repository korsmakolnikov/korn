package main

import (
	"log"

	"github.com/korsmakolnikov/kornvimgen/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
