package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/utils"
)

// SQLAddColumn return the add column sql for table create
func (grammarSQL SQL) SQLAddColumn(db *sqlx.DB, Column *grammar.Column) string {
	types := grammarSQL.Types
	quoter := grammarSQL.Quoter

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

	unsigned := utils.GetIF(Column.IsUnsigned, "UNSIGNED", "").(string)
	nullable := utils.GetIF(Column.Nullable, "NULL", "NOT NULL").(string)
	defaultValue := utils.GetIF(Column.Default != nil, fmt.Sprintf("DEFAULT %v", Column.Default), "").(string)
	comment := utils.GetIF(utils.StringVal(Column.Comment) != "", fmt.Sprintf("COMMENT %s", quoter.VAL(Column.Comment, db)), "").(string)
	collation := utils.GetIF(utils.StringVal(Column.Collation) != "", fmt.Sprintf("COLLATE %s", utils.StringVal(Column.Collation)), "").(string)
	extra := utils.GetIF(utils.StringVal(Column.Extra) != "", "AUTO_INCREMENT", "")
	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s %s",
		quoter.ID(Column.Name, db), typ, unsigned, nullable, defaultValue, extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

// SQLAddIndex  return the add index sql for table create
func (grammarSQL SQL) SQLAddIndex(db *sqlx.DB, index *grammar.Index) string {
	indexTypes := grammarSQL.IndexTypes
	quoter := grammarSQL.Quoter

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

	sql := fmt.Sprintf(
		"%s %s (%s) %s",
		typ, quoter.ID(index.Name, db), strings.Join(columns, ","), comment)

	return sql
}

// SQLAddPrimary return the add primary key sql for table create
func (grammarSQL SQL) SQLAddPrimary(db *sqlx.DB, primary *grammar.Primary) string {

	quoter := grammarSQL.Quoter

	// PRIMARY KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range primary.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	sql := fmt.Sprintf(
		"PRIMARY KEY %s (%s)",
		quoter.ID(primary.Name, db), strings.Join(columns, ","))

	return sql
}
