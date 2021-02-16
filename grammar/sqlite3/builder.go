package sqlite3

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/utils"
)

// SQLTableExists return the SQL for checking table exists.
func (builder Builder) SQLTableExists(db *sqlx.DB, name string, quoter grammar.Quoter) string {
	return fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name=%s", quoter.VAL(name, db))
}

// SQLAddIndex  return the add index sql for table create
func (builder Builder) SQLAddIndex(db *sqlx.DB, index *grammar.Index, indexTypes map[string]string, quoter grammar.Quoter) string {
	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	if typ == "KEY" {
		return ""
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range index.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, quoter.ID(index.Name, db), quoter.ID(index.TableName, db), strings.Join(columns, ","))

	return sql
}

// SQLRenameTable return the SQL for the renaming table.
func (builder Builder) SQLRenameTable(db *sqlx.DB, old string, new string, quoter grammar.Quoter) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", quoter.ID(old, db), quoter.ID(new, db))
}

// SQLAddColumn return the add column sql for table create
func (builder Builder) SQLAddColumn(db *sqlx.DB, Column *grammar.Column, types map[string]string, quoter grammar.Quoter) string {
	// `id` bigint(20) unsigned NOT NULL,
	typ, has := types[Column.Type]
	if !has {
		typ = "VARCHAR"
	}
	if Column.Precision != nil && Column.Scale != nil {
		typ = fmt.Sprintf("%s(%d,%d)", typ, utils.IntVal(Column.Precision), utils.IntVal(Column.Scale))
	} else if Column.DatetimePrecision != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(Column.DatetimePrecision))
	} else if Column.Length != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(Column.Length))
	}
	defaultValue := utils.GetIF(Column.Default != nil, fmt.Sprintf("DEFAULT %v", Column.Default), "").(string)
	unsigned := utils.GetIF(Column.IsUnsigned, "UNSIGNED", "").(string)
	primaryKey := utils.GetIF(Column.Primary, "PRIMARY KEY", "").(string)
	nullable := utils.GetIF(Column.Nullable, "NULL", "NOT NULL").(string)
	if defaultValue == "" && nullable == "NOT NULL" {
		nullable = "NULL"
	}

	if primaryKey != "" {
		nullable = primaryKey
	}
	comment := utils.GetIF(Column.Comment != nil, fmt.Sprintf("COMMENT %s", quoter.VAL(Column.Comment, db)), "").(string)
	collation := utils.GetIF(Column.Collation != nil, fmt.Sprintf("COLLATE %s", utils.StringVal(Column.Collation)), "").(string)
	extra := utils.GetIF(Column.Extra != nil, "AUTOINCREMENT", "")
	if extra == "AUTOINCREMENT" {
		unsigned = ""
	}

	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s %s",
		quoter.ID(Column.Name, db), typ, unsigned, nullable, defaultValue, extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
}
