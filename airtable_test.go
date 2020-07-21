package airtable

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAirtable(t *testing.T) {
	acc := Account{
		ApiKey: "abc",
		BaseId: "def",
	}
	tbl := NewTable("test", acc)
	Convey("NewTable returns GenericTable struct", t, func() {
		So(tbl, ShouldResemble, &GenericTable{
			account: &acc,
			Name:    "test",
			View:    VIEWNAME,
		})
	})

	Convey("List returns records", t, func() {
		ret, err := tbl.List(Options{})
		So(err, ShouldEqual, nil)
		So(ret, ShouldResemble, []Record{})
	})

	Convey("Update returns nil if no error", t, func() {
		err := tbl.Update([]Record{})
		So(err, ShouldEqual, nil)
	})

}