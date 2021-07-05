```
package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jordancurve/dbslice"
	_ "github.com/mattn/go-sqlite3"
)

type Place struct {
	Country       string `db:"country"`
	City          string `db:"city"`
	TelephoneCode int    `db:"telephone_code" dbslice:"primarykey"`
}

func main() {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("sqlx.Open() failed: %s", err))
	}
	want := []Place{
		{"usa", "sf", 7071234567},
		{"canada", "vancouver", 12345},
	}
	db.MustExec(dbslice.CreateTableSQL("places", []Place{}))
	dbslice.MustInsertSlice(db, "places", want)
	got := []Place{}
	dbslice.MustAppendToSlice(db, &got, "select * from places")
	fmt.Printf("got=%v\n", got)
}
```
