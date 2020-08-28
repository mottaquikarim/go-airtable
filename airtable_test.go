package airtable

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAirtable(t *testing.T) {
	acc := func() Account {
		return Account{
			ApiKey: "abc",
			BaseId: "def",
		}
	}

	testAcc := acc()
	tbl := NewTable("test", testAcc)
	Convey("NewTable returns GenericTable struct", t, func() {
		testAcc := acc()
		testAcc.BaseUrl = BASEURL
		So(tbl, ShouldResemble, &GenericTable{
			account: &testAcc,
			Name:    "test",
			View:    VIEWNAME,
		})
	})

	Convey("List returns records", t, func() {
		ret, err := tbl.List(Options{})
		So(err, ShouldEqual, nil)
		So(ret, ShouldResemble, []Record{})
	})

	Convey("a list of records is built through pagination", t, func(c C) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			is_first_request := len(r.URL.Query()["offset"]) == 0

			if is_first_request {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"records":[{"id":"testId","fields":{"HasRun":true}}],"offset":"basetoken/RandomOffset"}`)

			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"records":[{"id":"testId2","fields":{"HasRun":true}}]}`)
			}
		}

		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		records, err := tbl.List(Options{})

		So(err, ShouldEqual, nil)
		So(len(records), ShouldEqual, 2)
		So(records[0].ID, ShouldEqual, "testId")
		So(records[1].ID, ShouldEqual, "testId2")
	})

	Convey("an offset returned by the server is used in the following request", t, func(c C) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			is_first_request := len(r.URL.Query()["offset"]) == 0

			if is_first_request {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"records":[], "offset":"basetoken/RandomOffset"}`)

				c.So(r.URL.Query()["offset"], ShouldResemble, []string(nil))
			} else {
				c.So(r.URL.Query()["offset"], ShouldResemble, []string{"basetoken/RandomOffset"})
			}
		}

		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		_, _ = tbl.List(Options{})
	})

	Convey("offset is passed along as expected, if there is one", t, func(c C) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.So(r.URL.Query()["offset"], ShouldResemble, []string{"basetoken/RandomOffset"})
			_, _ = w.Write([]byte("test"))
		}))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		_, _ = tbl.List(Options{
			Offset: "basetoken/RandomOffset",
		})
	})

	Convey("max records is passed along as expected", t, func(c C) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.So(r.URL.Query()["maxRecords"], ShouldResemble, []string{"2"})
			_, _ = w.Write([]byte("test"))
		}))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		_, _ = tbl.List(Options{
			MaxRecords: 2,
		})
	})

	Convey("all filters work", t, func(c C) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.So(r.URL.Query()["maxRecords"], ShouldResemble, []string{fmt.Sprint(MAXRECORDS)})
			// no view explicitly set, so sends to table's "default view"
			c.So(r.URL.Query()["view"], ShouldResemble, []string{VIEWNAME})
			c.So(r.URL.Query()["filterByFormula"], ShouldResemble, []string{"NOT({HasRun})"})
			c.So(r.URL.Query()["sort[0][field]"], ShouldResemble, []string{"Date"})
			c.So(r.URL.Query()["sort[0][direction]"], ShouldResemble, []string{"desc"})
			_, _ = w.Write([]byte("test"))
		}))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		_, _ = tbl.List(Options{
			Filter: "NOT({HasRun})",
			Sort: []map[string]string{
				map[string]string{
					"field":     "Date",
					"direction": "desc",
				},
			},
		})
	})

	Convey("View is overrideable", t, func(c C) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.So(r.URL.Query()["view"], ShouldResemble, []string{"foobar"})
			_, _ = w.Write([]byte("test"))
		}))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL
		tbl := NewTable("hello", testAcc)
		_, _ = tbl.List(Options{
			View: "foobar",
		})
	})

	Convey("Update sends correct request body", t, func(c C) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			c.So(string(body), ShouldEqual, `{"records":[{"id":"testId","fields":{"HasRun":true}}]}`)
			_, _ = w.Write([]byte("test"))
		}))
		defer ts.Close()

		testAcc := acc()
		testAcc.BaseUrl = ts.URL

		tbl := NewTable("hello", testAcc)

		err := tbl.Update([]Record{
			Record{
				ID: "testId",
				Fields: map[string]interface{}{
					"HasRun": true,
				},
			},
		})
		So(err, ShouldEqual, nil)
	})

}
