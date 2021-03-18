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
	for _, column := range index.Columns {
		columns = append(columns, quoter.ID(column.Name, db))
	}

	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, quoter.ID(index.Name, db), quoter.ID(index.TableName, db), strings.Join(columns, ","))

	return sql
}

// SQLAddColumn return the add column sql for table create
func (grammarSQL SQLite3) SQLAddColumn(db *sqlx.DB, column *dbal.Column) string {

	quoter := grammarSQL.Quoter
	types := grammarSQL.Types

	// `id` bigint(20) unsigned NOT NULL,
	typ, has := types[column.Type]
	if !has {
		typ = "VARCHAR"
	}
	if column.Precision != nil && column.Scale != nil {
		typ = fmt.Sprintf("%s(%d,%d)", typ, utils.IntVal(column.Precision), utils.IntVal(column.Scale))
	} else if column.DateTimePrecision != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(column.DateTimePrecision))
	} else if typ == "BLOB" {
		typ = "BLOB"
	} else if typ == "ENUM" {
		option := fmt.Sprintf("('%s')", strings.Join(column.Option, "','"))
		typ = fmt.Sprintf("TEXT CHECK( %s IN %s )", quoter.ID(column.Name, db), option)
	} else if column.Length != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(column.Length))
	}
	defaultValue := utils.GetIF(column.Default != nil, fmt.Sprintf("DEFAULT %v", column.Default), "").(string)
	// unsigned := utils.GetIF(column.IsUnsigned && column.Type == "BIGINT", "UNSIGNED", "").(string)
	primaryKey := utils.GetIF(column.Primary, "PRIMARY KEY", "").(string)
	nullable := utils.GetIF(column.Nullable, "NULL", "NOT NULL").(string)
	if defaultValue == "" && nullable == "NOT NULL" {
		nullable = "NULL"
	}

	if primaryKey != "" {
		nullable = primaryKey
	}

	collation := utils.GetIF(column.Collation != nil, fmt.Sprintf("COLLATE %s", utils.StringVal(column.Collation)), "").(string)
	extra := utils.GetIF(column.Extra != nil, "AUTOINCREMENT", "")

	if extra == "AUTOINCREMENT" {
		typ = "INTEGER"
	}

	if column.IsUnsigned && typ == "BIGINT" {
		typ = "UNSIGNED BIG INT"
	}

	// JSON type
	if typ == "JSON" || typ == "JSONB" {
		typ = "TEXT"
	} else if typ == "UUID" { // uuid
		typ = "VARCHAR(36)"
	} else if typ == "IPADDRESS" { // ipAdderss
		typ = "integer"
	}

	sql := fmt.Sprintf(
		"%s %s %s %s %s %s",
		quoter.ID(column.Name, db), typ, nullable, defaultValue, extra, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

// SQLAddPrimary return the add primary key sql for table create
func (grammarSQL SQLite3) SQLAddPrimary(db *sqlx.DB, primary *dbal.Primary) string {

	quoter := grammarSQL.Quoter

	// PRIMARY KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, column := range primary.Columns {
		columns = append(columns, quoter.ID(column.Name, db))
	}

	sql := fmt.Sprintf(
		"PRIMARY KEY (%s)",
		strings.Join(columns, ","))

	return sql
}
