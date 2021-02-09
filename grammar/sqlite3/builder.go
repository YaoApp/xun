package sqlite3

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
)

// SQLTableExists return the SQL for checking table exists.
func (builder Builder) SQLTableExists(db *sqlx.DB, name string, quoter grammar.Quoter) string {
	return fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name=%s", quoter.VAL(name, db))
}

// SQLCreateIndex  return the add index sql for table create
func (builder Builder) SQLCreateIndex(db *sqlx.DB, index *grammar.Index, indexTypes map[string]string, quoter grammar.Quoter) string {
	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range index.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, quoter.ID(index.Name, db), quoter.ID(index.TableName, db), strings.Join(columns, "`,`"))

	return sql
}

// SQLRenameTable return the SQL for the renaming table.
func (builder Builder) SQLRenameTable(db *sqlx.DB, old string, new string, quoter grammar.Quoter) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", quoter.ID(old, db), quoter.ID(new, db))
}
