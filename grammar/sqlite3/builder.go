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
	for _, field := range index.Fields {
		columns = append(columns, quoter.ID(field.Field, db))
	}

	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, quoter.ID(index.Index, db), quoter.ID(index.TableName, db), strings.Join(columns, "`,`"))

	return sql
}
