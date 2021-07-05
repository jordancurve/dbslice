package dbslice

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type Place struct {
	Country       string `db:"country"`
	City          string `db:"city"`
	TelephoneCode int    `db:"telephone_code" dbslice:"primarykey"`
}

type X struct {
	x int `db:"x"`
}

func TestCreateTableSQL(t *testing.T) {
	tests := []struct {
		tableName string
		slice     interface{}
		opts      []Option
		want      string
	}{
		{"single_column", []X{}, nil, "CREATE TABLE single_column (x);"},
		{"single_column", []X{}, []Option{WithoutRowid()},
			"CREATE TABLE single_column (x) WITHOUT ROWID;"},
		{"single_column", []X{}, []Option{IfNotExists()},
			"CREATE TABLE IF NOT EXISTS single_column (x);"},
		{"single_column", []X{}, []Option{IfNotExists(), WithoutRowid()},
			"CREATE TABLE IF NOT EXISTS single_column (x) WITHOUT ROWID;"},
		{"places", []Place{}, nil,
			"CREATE TABLE places (country, city, telephone_code, PRIMARY KEY (telephone_code));"},
	}
	for _, test := range tests {
		got := CreateTableSQL(test.tableName, test.slice, test.opts...)
		if got != test.want {
			t.Errorf("CreateTable(%q, %v, %v)=%q; want %q",
				test.tableName, test.slice, test.opts, got, test.want)
		}
	}
}

func TestEverything(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sqlx.Open() failded: %s", err)
	}
	want := []Place{
		{"usa", "sf", 7071234567},
		{"canada", "vancouver", 12345},
	}
	db.MustExec(CreateTableSQL("places", []Place{}))
	MustInsertSlice(db, "places", want)
	got := []Place{}
	MustAppendToSlice(db, &got, "select * from places")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("test everything failed (-want, +got):\n%s", diff)
	}
	MustAppendToSlice(db, &got, "select * from places")
	if diff := cmp.Diff(append(want, want...), got); diff != "" {
		t.Errorf("test everything failed (-want, +got):\n%s", diff)
	}
}
