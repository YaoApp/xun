package sqlite3

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
)

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
