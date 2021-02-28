package schema

import (
	"fmt"

	"github.com/yaoapp/xun/utils"
)

// NewColumn Create a new column instance
func (table *Table) NewColumn(name string) *Column {
	column := &Column{
		Column: table.Table.NewColumn(name),
		Table:  table,
	}
	return column
}

// PushColumn add a column to the table
func (table *Table) PushColumn(column *Column) *Table {
	table.Table.PushColumn(column.Column)
	table.ColumnMap[column.Name] = column
	return table
}

// GetColumn get the column instance of the table, if the column does not exist return nil.
func (table *Table) GetColumn(name string) *Column {
	return table.ColumnMap[name]
}

// Column get the column instance of the table, if the column does not exist create a new one.
func (table *Table) Column(name string) *Column {
	column, has := table.ColumnMap[name]
	if has {
		return column
	}
	return table.NewColumn(name)
}

// HasColumn Determine if the table has a given column.
func (table *Table) HasColumn(name ...string) bool {
	has := true
	for _, n := range name {
		_, has = table.ColumnMap[n]
		if !has {
			return has
		}
	}
	return has
}

// PutColumn add or modify a column to the table
func (table *Table) PutColumn(column *Column) *Table {
	if table.HasColumn(column.Name) {
		table.ChangeColumn(column)
	} else {
		table.AddColumn(column)
	}
	return table
}

// AddColumn add a column to the table
func (table *Table) AddColumn(column *Column) *Table {
	table.PushColumn(column)
	table.AddColumnCommand(column.Column, nil, func() {
		delete(table.ColumnMap, column.Name)
	})
	return table
}

// ChangeColumn modify a column to the table
func (table *Table) ChangeColumn(column *Column) *Table {
	table.ChangeColumnCommand(column.Column, func() {
		table.ColumnMap[column.Name] = column
	}, nil)
	return table
}

// DropColumn Indicate that the given columns should be dropped.
func (table *Table) DropColumn(name ...string) {
	for _, n := range name {
		table.DropColumnCommand(n, func() {
			delete(table.ColumnMap, n)
		}, nil)
	}
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Table) RenameColumn(old string, new string) *Column {
	column := table.GetColumn(old)
	column.Name = new
	table.ColumnMap[new] = column
	table.RenameColumnCommand(old, new, func() {
		delete(table.ColumnMap, old)
	}, func() {
		delete(table.ColumnMap, new)
	})
	return column
}

// HasIndex check if the column has created the index
func (column *Column) HasIndex(name string) bool {
	for _, idx := range column.Indexes {
		if idx.Name == name {
			return true
		}
	}
	return false
}

// Unique set as index
func (column *Column) Unique() *Column {
	name := fmt.Sprintf("%s_%s", column.Name, "unique")
	if column.HasIndex(name) {
		return column
	}
	column.Table.AddUnique(name, column.Name)
	return column
}

// Primary set as primary key
func (column *Column) Primary() *Column {
	if column.Column.Primary {
		return column
	}
	column.Column.Primary = true
	column.Table.AddPrimary(column.Name)
	return column
}

// Index set as index key
func (column *Column) Index() *Column {
	name := fmt.Sprintf("%s_%s", column.Name, "index")
	if column.HasIndex(name) {
		return column
	}
	column.Table.AddIndex(name, column.Name)
	return column
}

// Unsigned set the column IsUnsigned attribute is true
func (column *Column) Unsigned() *Column {
	column.IsUnsigned = true
	return column
}

// Null set the column nullable attribute is true
func (column *Column) Null() *Column {
	column.Nullable = true
	return column
}

// NotNull set the column nullable attribute is false
func (column *Column) NotNull() *Column {
	column.Nullable = false
	return column
}

// AutoIncrement set the numeric column AutoIncrement attribute is true
func (column *Column) AutoIncrement() *Column {
	column.Extra = utils.StringPtr("AutoIncrement")
	return column
}
