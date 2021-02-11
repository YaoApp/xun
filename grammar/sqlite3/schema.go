package sqlite3

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/logger"
)

// Create a new table on the schema
func (grammar SQLite3) Create(table *grammar.Table, db *sqlx.DB) error {
	name := grammar.Quoter.ID(table.Name, db)
	sql := fmt.Sprintf("CREATE TABLE %s (\n", name)
	stmts := []string{}

	// Columns
	for _, Column := range table.Columns {
		stmts = append(stmts,
			grammar.Builder.SQLCreateColumn(db, Column, grammar.Types, grammar.Quoter),
		)
	}
	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + fmt.Sprintf("\n)")

	// Create table
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
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
	defer logger.Debug(logger.CREATE, indexes...).TimeCost(time.Now())
	_, err = db.Exec(strings.Join(indexes, ";\n"))
	if err != nil {
		return err
	}

	return nil
}
