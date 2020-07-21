package main

import (
	"fmt"
	"os"

	airtable "github.com/mottaquikarim/go-airtable"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

var (
	fs     = flag.NewFlagSetWithEnvPrefix(os.Args[0], "AIRTABLE", 0)
	apiKey = fs.String("api-key", "XXXXX", "Airtable API Key")
	baseId = fs.String("base-id", "XXXXX", "Airtable Base Id")
)

func usage() {
	fmt.Println("Usage: ./airtbl [flags]")
	fs.PrintDefaults()
}

func main() {
	fs.Usage = usage
	fs.Parse(os.Args[1:])

	account := airtable.Account{
		ApiKey: *apiKey,
		BaseId: *baseId,
	}

	pokédex := airtable.NewTable("pokémon", account)
	original_generation, err := pokédex.List(airtable.Options{
		MaxRecords: 100,
		View:       "All",
	})
	if err != nil {
		// handle error
		log.Printf("Error! %v", err)
	}

	for _, pokémon := range original_generation {
		log.Printf("ID: %s Name: %v", pokémon.Fields["ID"], pokémon.Fields["Name"])
	}
}
