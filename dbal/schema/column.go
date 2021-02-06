package schema

import (
	"errors"
	"fmt"
	"strings"
)

var columnTypes = map[string]string{
	"string": "VARCHAR",
}

var fieldTypes = map[string]string{
	"VARCHAR": "string",
}

// Unique set as index
func (column *Column) Unique() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "unique"
	column.Table.addIndex(index)
	return column
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Blueprint) BigInteger() {}

// String Create a new string column on the table.
func (table *Blueprint) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = &length
	column.Type = "string"
	table.addColumn(column)
	return column
}

// Change mark as changed
func (column *Column) Change() *Column {
	column.Changed = true
	return column
}

// Remove mark as removed
func (column *Column) Remove() *Column {
	column.Removed = true
	return column
}

// Rename rename the column
func (column *Column) Rename(name string) *Column {
	column.Changed = true
	column.Newname = name
	return column
}

// UpField update the column by given table field.
func (column *Column) UpField(field *TableField) *Column {
	column.Name = field.Field
	column.Type = GetColumnType(field.Type)
	return column
}

func (column *Column) sqlCreate() string {
	// `id` bigint(20) unsigned NOT NULL,
	sql := fmt.Sprintf("`%s` %s(%d) %s", column.nameEscaped(), GetTableFieldType(column.Type), *column.Length, "NOT NULL")
	return sql
}

func (column *Column) validate() *Column {
	if column.nameEscaped() == "" {
		err := errors.New("the column name must be set")
		panic(err)
	}

	if column.Type == "" {
		err := errors.New("the column " + column.Name + " type must be set")
		panic(err)
	}

	if column.Table == nil {
		err := errors.New("the column " + column.Name + "does not bind the table")
		panic(err)
	}
	return column
}

// GetColumnType return the columns type
func GetColumnType(name string) string {
	if _, has := fieldTypes[name]; has {
		return fieldTypes[name]
	}
	return "varchar"
}

// GetTableFieldType return the columns type
func GetTableFieldType(name string) string {
	if _, has := columnTypes[name]; has {
		return columnTypes[name]
	}
	return "string"
}

func (column *Column) nameEscaped() string {
	return strings.ReplaceAll(column.Name, "`", "")
}
