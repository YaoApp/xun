package sqlite3

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// SQLAddIndex  return the add index sql for table create
func (grammarSQL SQLite3) SQLAddIndex(db *sqlx.DB, index *dbal.Index) string {
	quoter := grammarSQL.Quoter
	indexTypes := grammarSQL.IndexTypes
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

// SQLAddColumn return the add column sql for table create
func (grammarSQL SQLite3) SQLAddColumn(db *sqlx.DB, Column *dbal.Column) string {

	quoter := grammarSQL.Quoter
	types := grammarSQL.Types

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

// SQLAddPrimary return the add primary key sql for table create
func (grammarSQL SQLite3) SQLAddPrimary(db *sqlx.DB, primary *dbal.Primary) string {

	quoter := grammarSQL.Quoter

	// PRIMARY KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range primary.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	sql := fmt.Sprintf(
		"PRIMARY KEY (%s)",
		strings.Join(columns, ","))

	return sql
}
