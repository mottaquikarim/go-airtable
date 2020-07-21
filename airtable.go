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

func NewTable(name string, account Account) Table {
	return &GenericTable{
		account: &account,
		Name:    name,
		View:    VIEWNAME,
	}
}

func (t *GenericTable) List(opts Options) ([]Record, error) {
	log.Printf("calling list")
	return []Record{}, nil
}

func (t *GenericTable) Update(recs []Record) error {
	log.Printf("calling update")
	return nil
}
