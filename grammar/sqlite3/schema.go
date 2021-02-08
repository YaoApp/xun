package sqlite3

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
)

// Exists the Exists
func (grammar SQLite3) Exists(name string, db *sqlx.DB) bool {
	sql := fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name=%s", grammar.Quoter.VAL(name, db))
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
func (grammar SQLite3) Create(table *grammar.Table, db *sqlx.DB) error {
	name := grammar.Quoter.ID(table.Name, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}

	// fields
	for _, field := range table.Fields {
		stmts = append(stmts,
			grammar.Builder.SQLCreateColumn(db, field, grammar.Types, grammar.Quoter),
		)
	}
	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf("\n)")

	// Create table
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	// indexes
	indexes := []string{}
	for _, index := range table.Indexes {
		indexes = append(indexes,
			grammar.Builder.SQLCreateIndex(db, index, grammar.IndexTypes, grammar.Quoter),
		)
	}

	_, err = db.Exec(strings.Join(indexes, ";\n"))
	if err != nil {
		return err
	}

	return nil
}
