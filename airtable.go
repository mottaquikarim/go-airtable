// Package airtable provides a simple client for interacting
// with the Airtable API.
package airtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// getFullUrl builds the URL for making API call
func (t *GenericTable) getFullUrl() string {
	return fmt.Sprintf("%s/v%s/%s/%s", t.account.BaseUrl, VERSION, t.account.BaseId, t.Name)
}

// NewTable returns a GenericTable struct that implements
// the Table interface
func NewTable(name string, account Account) Table {
	if len(account.BaseUrl) == 0 {
		account.BaseUrl = BASEURL
	}

	return &GenericTable{
		account: &account,
		Name:    name,
		View:    VIEWNAME,
	}
}

// List returns a list of records from the Airtable.
func (t *GenericTable) List(opts Options) ([]Record, error) {
	// create req
	req, err := http.NewRequest("GET", t.getFullUrl(), nil)
	if err != nil {
		return []Record{}, fmt.Errorf("Failed to create request object")
	}

	// init query params
	q := req.URL.Query()

	// add maxrecords
	switch opts.MaxRecords {
	case 0:
		q.Add("maxRecords", fmt.Sprint(MAXRECORDS))
	default:
		q.Add("maxRecords", fmt.Sprint(opts.MaxRecords))
	}

	// add view
	switch len(opts.View) {
	case 0:
		q.Add("view", fmt.Sprint(VIEWNAME))
	default:
		q.Add("view", fmt.Sprint(opts.View))
	}

	// add filters if they exist
	if len(opts.Filter) > 0 {
		q.Add("filterByFormula", opts.Filter)
	}

	// add sorting if provided
	// sorting must be assembled as follows:
	// [{"field": "my-field-name", "direction": "asc|desc"}]
	// this is converted to:
	// 		sort[0][field]=my-field-name
	// 		sort[0][direction]=asc|desc
	if len(opts.Sort) > 0 {
		for i, sort := range opts.Sort {
			for key, val := range sort {
				q.Add(fmt.Sprintf("sort[%d][%s]", i, key), val)
			}
		}
	}

	// encode everything
	req.URL.RawQuery = q.Encode()

	// set headsers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.account.ApiKey))

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Record{}, fmt.Errorf("Error occured: %v", err)
	}
	defer resp.Body.Close()

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Record{}, fmt.Errorf("Failed to read body: %v", err)
	}

	// unmarshal
	ret := records{}
	if err = json.Unmarshal(body, &ret); err != nil {
		return []Record{}, fmt.Errorf("Failed to unmarshal response: %v", err)
	}

	return ret.Records, nil
}

// Update makes a PATCH request to all records provided to Airtable.
func (t *GenericTable) Update(recs []Record) error {
	// assemble body
	recordWrapper := records{
		Records: recs,
	}
	jsonStr, err := json.Marshal(recordWrapper)
	if err != nil {
		return fmt.Errorf("Failed to create request body %v", err)
	}

	// create req
	req, err := http.NewRequest("PATCH", t.getFullUrl(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("Failed to create new request %v", err)
	}

	// set headsers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.account.ApiKey))
	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to make request %v", err)
	}
	defer resp.Body.Close()

	return nil
}
