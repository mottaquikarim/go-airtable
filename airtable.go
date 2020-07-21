// Package airtable provides a simple client for interacting
// with the Airtable API.
package airtable

import (
	// "encoding/json"
	// "fmt"
	// "net/http"

	log "github.com/sirupsen/logrus"
)

type Table interface {
	List(opts Options) ([]Record, error)
	Update(recs []Record) error
}

type GenericTable struct {
	account *Account
	Name    string
	View    string
}

// NewTable returns a GenericTable struct that implements
// the Table interface
func NewTable(name string, account Account) Table {
	return &GenericTable{
		account: &account,
		Name:    name,
		View:    VIEWNAME,
	}
}

// List returns a list of records from the Airtable.
func (t *GenericTable) List(opts Options) ([]Record, error) {
	log.Printf("calling list")
	return []Record{}, nil
}

// Update makes a PATCH request to all records provided to Airtable.
func (t *GenericTable) Update(recs []Record) error {
	log.Printf("calling update")
	return nil
}
