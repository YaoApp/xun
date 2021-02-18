package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/utils"
)

// SQLTableExists return the SQL for checking table exists.
func (builder Builder) SQLTableExists(db *sqlx.DB, name string, quoter grammar.Quoter) string {
	return fmt.Sprintf("SELECT table_name AS name FROM information_schema.tables WHERE table_name = %s", quoter.VAL(name, db))
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

	unsigned := ""
	if Column.IsUnsigned {
		typ = "BIGINT"
	}
	nullable := utils.GetIF(Column.Nullable, "NULL", "NOT NULL").(string)
	defaultValue := utils.GetIF(Column.Default != nil, fmt.Sprintf("DEFAULT %v", Column.Default), "").(string)
	comment := utils.GetIF(utils.StringVal(Column.Comment) != "", fmt.Sprintf("COMMENT %s", quoter.VAL(Column.Comment, db)), "").(string)
	collation := utils.GetIF(utils.StringVal(Column.Collation) != "", fmt.Sprintf("COLLATE %s", utils.StringVal(Column.Collation)), "").(string)
	extra := ""
	if utils.StringVal(Column.Extra) != "" {
		typ = "SERIAL"
		nullable = ""
		defaultValue = ""
	}
	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s %s",
		quoter.ID(Column.Name, db), typ, unsigned, nullable, defaultValue, extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

// SQLAddIndex  return the add index sql for table create
func (builder Builder) SQLAddIndex(db *sqlx.DB, index *grammar.Index, indexTypes map[string]string, quoter grammar.Quoter) string {
	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range index.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	comment := ""
	if index.Comment != nil {
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(index.Comment, db))
	}
	name := quoter.ID(index.Name, db)
	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, name, quoter.ID(index.TableName, db), strings.Join(columns, ","))

	if typ == "PRIMARY KEY" {
		sql = fmt.Sprintf(
			"%s (%s) %s",
			typ, strings.Join(columns, ","), comment)
	}
	return sql
}
