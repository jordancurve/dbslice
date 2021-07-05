```
package main

import (
	"github/jordancurve/dbslice"
	"github.com/jmoiron/sqlx"
)

type Place struct {
	Country       string `db:"country"`
	City          string `db:"city"`
	TelephoneCode int    `db:"telephone_code" dbslice:"primarykey"`
}

func main() {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sqlx.Open() failded: %s", err)
	}
	want := []Place{
		{"usa", "sf", 7071234567},
		{"canada", "vancouver", 12345},
	}
	db.MustExec(CreateTableSQL("places", []Place{}))
	dbslice.MustInsertSlice(db, "places", want)
	got := []Place{}
	dbslice.MustAppendToSlice(db, &got, "select * from places")
}
```
