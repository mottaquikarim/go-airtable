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

	log.Printf("Here %v", account)
}
