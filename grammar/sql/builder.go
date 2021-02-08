package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
)

// Builder the SQL builder
type Builder struct{}

// SQLCreateColumn return the add column sql for table create
func (builder Builder) SQLCreateColumn(db *sqlx.DB, field *grammar.Field, types map[string]string, quoter grammar.Quoter) string {
	// `id` bigint(20) unsigned NOT NULL,
	typ, has := types[field.Type]
	if !has {
		typ = "VARCHAR"
	}
	if field.Precision > 0 && field.Scale > 0 {
		typ = fmt.Sprintf("%s(%d,%d)", typ, field.Precision, field.Scale)
	} else if field.DatetimePrecision > 0 {
		typ = fmt.Sprintf("%s(%d)", typ, field.DatetimePrecision)
	} else if field.Length > 0 {
		typ = fmt.Sprintf("%s(%d)", typ, field.Length)
	}

	nullable := "NOT NULL"
	if field.Nullable {
		nullable = " NULL"
	}

	defaultValue := ""
	if field.Default != nil {
		defaultValue = fmt.Sprintf("DEFAULT %v", field.Default)
	}

	comment := ""
	if field.Comment != "" {
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(field.Comment, db))
	}

	collation := ""
	if field.Collation != "" {
		collation = fmt.Sprintf("COLLATE %s", collation)
	}

	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s",
		quoter.ID(field.Field, db), typ, nullable, defaultValue, field.Extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
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

	comment := ""
	if index.Comment != "" {
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(index.Comment, db))
	}

	sql := fmt.Sprintf(
		"%s %s (%s) %s",
		typ, quoter.ID(index.Index, db), strings.Join(columns, "`,`"), comment)

	return sql
}