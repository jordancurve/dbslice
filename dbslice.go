package dbslice

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type Option int

const (
	ifNotExists  = Option(1)
	withoutRowid = Option(2)
)

func IfNotExists() Option {
	return ifNotExists
}

func WithoutRowid() Option {
	return withoutRowid
}

func CreateTableSQL(tableName string, slice interface{}, opts ...Option) string {
	sqlStmt := "CREATE TABLE "
	opt := map[Option]bool{}
	for _, o := range opts {
		opt[o] = true
	}
	if opt[ifNotExists] {
		sqlStmt += "IF NOT EXISTS "
	}
	sqlStmt += tableName + " ("
	s := reflect.TypeOf(slice).Elem() // struct
	columns := []string{}
	primaryKeys := []string{}
	for i := 0; i < s.NumField(); i++ {
		column := s.Field(i).Tag.Get("db")
		if column == "" {
			panic(fmt.Sprintf("CreateTableSQL: missing column struct tag for field %q", s.Field(i).Name))
		}
		if s.Field(i).Tag.Get("dbslice") == "primarykey" {
			primaryKeys = append(primaryKeys, column)
		}
		columns = append(columns, column)
	}
	if len(primaryKeys) > 0 {
		columns = append(columns, "PRIMARY KEY ("+strings.Join(primaryKeys, ", ")+")")
	}
	sqlStmt += strings.Join(columns, ", ") + ")"
	if opt[withoutRowid] {
		sqlStmt += " WITHOUT ROWID"
	}
	sqlStmt += ";"
	return sqlStmt
}

func MustInsertSlice(db *sqlx.DB, tableName string, slice interface{}) {
	s := reflect.TypeOf(slice).Elem() // struct
	qs := []string{}
	for i := 0; i < s.NumField(); i++ {
		qs = append(qs, "?")
	}
	insertStmt := "INSERT INTO " + tableName + " VALUES (" + strings.Join(qs, ", ") + ");"
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	v := reflect.ValueOf(slice)
	for i := 0; i < v.Len(); i++ {
		s := reflect.ValueOf(v.Index(i).Interface())
		args := []interface{}{}
		for j := 0; j < s.NumField(); j++ {
			args = append(args, s.Field(j).Interface())
		}
		if _, err := tx.Exec(insertStmt, args...); err != nil {
			panic(err)
		}
	}
	tx.Commit()
}

func MustAppendToSlice(db *sqlx.DB, slice interface{}, query string, args ...interface{}) {
	if err := sqlx.Select(db, slice, query, args...); err != nil {
		panic(fmt.Sprintf("%q failed: %s", query, args))
	}
}
