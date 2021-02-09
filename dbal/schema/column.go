package schema

import (
	"errors"
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
func (table *Table) BigInteger() {}

// String Create a new string column on the table.
func (table *Table) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = length
	column.Type = "string"
	table.addColumn(column)
	return column
}

// Drop mark as dropped for the index
func (column *Column) Drop() {
	column.dropped = true
}

// Rename mark as renamed with the given name for the index
func (column *Column) Rename(new string) {
	column.renamed = true
	column.newname = new
}

// Precision get the column Precision
func (column *Column) Precision() int {
	switch column.Args.(type) {
	case []int:
		return column.Args.([]int)[0]
	case int:
		return column.Args.(int)
	default:
		return 0
	}
}

// Scale get the column scale
func (column *Column) Scale() int {
	switch column.Args.(type) {
	case []int:
		if len(column.Args.([]int)) >= 1 {
			return column.Args.([]int)[1]
		}
		return 0
	default:
		return 0
	}
}

// DatetimePrecision get the column datetime precision
func (column *Column) DatetimePrecision() int {
	if column.Type != "timestamp" && column.Type != "datetime" && column.Type != "time" {
		return 0
	}
	switch column.Args.(type) {
	case []int:
		return column.Args.([]int)[0]
	case int:
		return column.Args.(int)
	default:
		return 0
	}
}

func (column *Column) validate() *Column {
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
	return column
}
