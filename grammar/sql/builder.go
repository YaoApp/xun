package sql

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// SQLAddColumn return the add column sql for table create
func (grammarSQL SQL) SQLAddColumn(column *dbal.Column) string {
	db := grammarSQL.DB
	quoter := grammarSQL.Quoter

	// `id` bigint(20) unsigned NOT NULL,
	typ := grammarSQL.getType(column)

	unsigned := utils.GetIF(column.IsUnsigned, "UNSIGNED", "").(string)
	nullable := utils.GetIF(column.Nullable, "NULL", "NOT NULL").(string)
	defaultValue := utils.GetIF(column.Default != nil, fmt.Sprintf("DEFAULT %v", quoter.VAL(column.Default, db)), "").(string)
	comment := utils.GetIF(utils.StringVal(column.Comment) != "", fmt.Sprintf("COMMENT %s", quoter.VAL(column.Comment, db)), "").(string)
	collation := utils.GetIF(utils.StringVal(column.Collation) != "", fmt.Sprintf("COLLATE %s", utils.StringVal(column.Collation)), "").(string)
	extra := utils.GetIF(utils.StringVal(column.Extra) != "", "AUTO_INCREMENT", "")

	if nullable == "NOT NULL" && strings.Contains(typ, "TIMESTAMP") && defaultValue == "" {
		if column.DateTimePrecision != nil {
			defaultValue = fmt.Sprintf("DEFAULT CURRENT_TIMESTAMP(%d)", *column.DateTimePrecision)
		} else {
			defaultValue = "DEFAULT CURRENT_TIMESTAMP"
		}
	}

	// JSON type
	if typ == "JSON" || typ == "JSONB" {
		mysql5_7_8, _ := semver.Make("5.7.8")
		version, err := grammarSQL.GetVersion()
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(fmt.Sprintf("T:%s|%s", column.Type, utils.StringVal(column.Comment)), db))
		if err != nil || version.LT(mysql5_7_8) {
			typ = "TEXT"
		} else {
			typ = "JSON"
		}
	} else if typ == "UUID" { // UUID
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(fmt.Sprintf("T:%s|%s", column.Type, utils.StringVal(column.Comment)), db))
		typ = "VARCHAR(36)"
	} else if typ == "IPADDRESS" { // ipAddress
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(fmt.Sprintf("T:%s|%s", column.Type, utils.StringVal(column.Comment)), db))
		typ = "integer"
	} else if typ == "MACADDRESS" { // macAddress 08:00:2b:01:02:03:04:05  bigint unsigned (8 bytes)
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(fmt.Sprintf("T:%s|%s", column.Type, utils.StringVal(column.Comment)), db))
		typ = "BIGINT"
		unsigned = "UNSIGNED"
	} else if typ == "YEAR" { // 2021 -1046 smallInt (2-byte)
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(fmt.Sprintf("T:%s|%s", column.Type, utils.StringVal(column.Comment)), db))
		typ = "SMALLINT"
	}

	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s %s",
		quoter.ID(column.Name, db), typ, unsigned, nullable, defaultValue, extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

func (grammarSQL SQL) getType(column *dbal.Column) string {
	// `id` bigint(20) unsigned NOT NULL,
	typ, has := grammarSQL.Types[column.Type]
	if !has {
		typ = "VARCHAR"
	}

	decimalTypes := []string{"DECIMAL", "FLOAT", "DOUBLE"}
	if utils.StringHave(decimalTypes, typ) && column.Precision != nil && column.Scale != nil {
		typ = fmt.Sprintf("%s(%d,%d)", typ, utils.IntVal(column.Precision), utils.IntVal(column.Scale))
	} else if column.DateTimePrecision != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(column.DateTimePrecision))
	} else if typ == "ENUM" {
		typ = fmt.Sprintf("ENUM('%s')", strings.Join(column.Option, "','"))
	} else if column.Length != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(column.Length))
	}

	return typ
}

// SQLAddIndex  return the add index sql for table create
func (grammarSQL SQL) SQLAddIndex(index *dbal.Index) string {
	db := grammarSQL.DB
	maxKeyLength := 256
	indexTypes := grammarSQL.IndexTypes
	quoter := grammarSQL.Quoter

	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, column := range index.Columns {
		if column.Type == "text" || column.Type == "mediumText" || column.Type == "longText" {
			columns = append(columns, fmt.Sprintf("%s(%d)", quoter.ID(column.Name, db), maxKeyLength))
		} else if column.Type == "json" || column.Type == "jsonb" { // ignore json and jsonb
			continue
		} else {
			columns = append(columns, quoter.ID(column.Name, db))
		}
	}

	comment := ""
	if index.Comment != nil {
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(index.Comment, db))
	}

	if len(columns) == 0 {
		return ""
	}

	sql := fmt.Sprintf(
		"%s %s (%s) %s",
		typ, quoter.ID(index.Name, db), strings.Join(columns, ","), comment)

	return sql
}

// SQLAddPrimary return the add primary key sql for table create
func (grammarSQL SQL) SQLAddPrimary(primary *dbal.Primary) string {

	quoter := grammarSQL.Quoter

	// PRIMARY KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, column := range primary.Columns {
		columns = append(columns, quoter.ID(column.Name, grammarSQL.DB))
	}

	sql := fmt.Sprintf(
		// "PRIMARY KEY %s (%s)",
		"PRIMARY KEY(%s)",
		// quoter.ID(primary.Name, grammarSQL.DB),
		strings.Join(columns, ","))

	return sql
}
