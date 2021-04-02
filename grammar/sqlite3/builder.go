package sqlite3

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// SQLAddIndex  return the add index sql for table create
func (grammarSQL SQLite3) SQLAddIndex(index *dbal.Index) string {
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
		columns = append(columns, quoter.ID(column.Name))
	}

	name := fmt.Sprintf("%s_%s", index.TableName, index.Name)
	sql := fmt.Sprintf(
		"CREATE %s %s ON %s (%s)",
		typ, quoter.ID(name), quoter.ID(index.TableName), strings.Join(columns, ","))

	return sql
}

// SQLAddColumn return the add column sql for table create
func (grammarSQL SQLite3) SQLAddColumn(column *dbal.Column) string {
	quoter := grammarSQL.Quoter

	// `id` bigint(20) unsigned NOT NULL,
	typ := grammarSQL.getType(column)

	defaultValue := utils.GetIF(column.Default != nil, fmt.Sprintf("DEFAULT %v", quoter.VAL(column.Default)), "").(string)
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

	sql := fmt.Sprintf(
		"%s %s %s %s %s %s",
		quoter.ID(column.Name), typ, nullable, defaultValue, extra, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

// getType
func (grammarSQL SQLite3) getType(column *dbal.Column) string {

	// `id` bigint(20) unsigned NOT NULL,
	typ, has := grammarSQL.Types[column.Type]
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
		typ = fmt.Sprintf("TEXT CHECK( %s IN %s )", grammarSQL.ID(column.Name), option)
	} else if column.Length != nil {
		typ = fmt.Sprintf("%s(%d)", typ, utils.IntVal(column.Length))
	}

	switch typ {
	case "JSON", "JSONB":
		typ = "TEXT"
		break
	case "UUID":
		typ = "VARCHAR(36)"
		break
	case "IPADDRESS": // 192.168.0.3
		typ = "integer"
		break
	case "MACADDRESS":
		typ = "UNSIGNED BIG INT" // / macAddress 08:00:2b:01:02:03:04:05  bigint unsigned (8 bytes)
		break
	case "YEAR":
		typ = "SMALLINT" // 2021 -1046
		break
	}

	return typ
}

// SQLAddPrimary return the add primary key sql for table create
func (grammarSQL SQLite3) SQLAddPrimary(primary *dbal.Primary) string {
	quoter := grammarSQL.Quoter

	// PRIMARY KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, column := range primary.Columns {
		columns = append(columns, quoter.ID(column.Name))
	}

	sql := fmt.Sprintf(
		// "CONSTRAINT %s PRIMARY KEY (%s)",
		"PRIMARY KEY (%s)",
		// quoter.ID(primary.Table.GetName()+"_pkey"),
		strings.Join(columns, ","))

	return sql
}
