package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
)

// Exists the Exists
func (grammar SQL) Exists(name string, db *sqlx.DB) bool {
	sql := grammar.Builder.SQLTableExists(db, name, grammar.Quoter)
	row := db.QueryRowx(sql)
	if row.Err() != nil {
		panic(row.Err())
	}
	res, err := row.SliceScan()
	if err != nil {
		return false
	}
	return name == fmt.Sprintf("%s", res[0])
}

// Create a new table on the schema
func (grammar SQL) Create(table *grammar.Table, db *sqlx.DB) error {
	name := grammar.Quoter.ID(table.Name, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}

	// fields
	for _, field := range table.Fields {
		stmts = append(stmts,
			grammar.Builder.SQLCreateColumn(db, field, grammar.Types, grammar.Quoter),
		)
	}

	// indexes
	for _, index := range table.Indexes {
		stmts = append(stmts,
			grammar.Builder.SQLCreateIndex(db, index, grammar.IndexTypes, grammar.Quoter),
		)
	}

	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf(
		"\n) ENGINE=%s DEFAULT CHARSET=%s COLLATE=%s",
		table.Engine, table.Charset, table.Collation,
	)

	_, err := db.Exec(sql)
	return err
}

// Drop a table from the schema.
func (grammar SQL) Drop(name string, db *sqlx.DB) error {
	sql := fmt.Sprintf("DROP TABLE %s", grammar.Quoter.ID(name, db))
	_, err := db.Exec(sql)
	return err
}
