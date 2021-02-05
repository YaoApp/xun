package schema

import (
	"errors"
	"fmt"
)

var columnTypes = map[string]string{
	"string": "VARCHAR",
}

// Unique set as index
func (column *Column) Unique() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "unique"
	column.Table.addIndex(index)
	return column
}

func (column *Column) sqlCreate() string {
	// `id` bigint(20) unsigned NOT NULL,
	sql := fmt.Sprintf("`%s` %s(%d) %s", column.Name, GetColumnType(column.Type), *column.Length, "NOT NULL")
	return sql
}

func (column *Column) validate() {
	if column.Name == "" {
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
}

// GetColumnType return the columns type
func GetColumnType(name string) string {
	if _, has := columnTypes[name]; has {
		return columnTypes[name]
	}
	return "varchar"
}
