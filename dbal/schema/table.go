package schema

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// GetName get the table name
func (table *Table) GetName() string {
	return table.Name
}

// GetPrefix get the table prefix
func (table *Table) GetPrefix() string {
	return table.Prefix
}

// GetFullName get the table name with prefix
func (table *Table) GetFullName() string {
	return table.TableName
}

// NewTable create a new blueprint intance
func NewTable(name string, builder *Builder) *Table {
	tableName := fmt.Sprintf("%s%s", builder.Conn.Option.Prefix, name)
	table := &Table{
		Name:      name,
		Prefix:    builder.Conn.Option.Prefix,
		Table:     dbal.NewTable(tableName, builder.SchemaName, builder.DBName),
		Builder:   builder,
		ColumnMap: map[string]*Column{},
		IndexMap:  map[string]*Index{},
	}
	return table
}

// GetColumns Get the columns map of the table
func (table *Table) GetColumns() map[string]*Column {
	return table.ColumnMap
}

// GetIndexes Get the indexes map of the table
func (table *Table) GetIndexes() map[string]*Index {
	return table.IndexMap
}

// Get Get the DBAL table instance
func (table *Table) Get() *Table {
	return table
}

// AddCommand Add a new command to the table.
func (table *Table) AddCommand(name string, success func(), fail func(), params ...interface{}) {
	table.Table.AddCommand(name, success, fail, params...)
}

// AddColumnCommand add a new command that adding a column
func (table *Table) AddColumnCommand(column *dbal.Column, success func(), fail func()) {
	table.AddCommand("AddColumn", success, fail, column)
}

// ChangeColumnCommand add a new command that modifing a column
func (table *Table) ChangeColumnCommand(column *dbal.Column, success func(), fail func()) {
	table.AddCommand("ChangeColumn", success, fail, column)
}

// RenameColumnCommand add a new command that renaming a column
func (table *Table) RenameColumnCommand(old string, new string, success func(), fail func()) {
	table.AddCommand("RenameColumn", success, fail, old, new)
}

// DropColumnCommand add a new command that dropping a column
func (table *Table) DropColumnCommand(name string, success func(), fail func()) {
	table.AddCommand("DropColumn", success, fail, name)
}

// CreateIndexCommand add a new command that creating a index
func (table *Table) CreateIndexCommand(index *dbal.Index, success func(), fail func()) {
	table.AddCommand("CreateIndex", success, fail, index)
}

// CreatePrimaryCommand add a new command that creating the primary key
func (table *Table) CreatePrimaryCommand(primary *dbal.Primary, success func(), fail func()) {
	table.AddCommand("CreatePrimary", success, fail, primary)
}

// DropPrimaryCommand add a new command drop the primary key
func (table *Table) DropPrimaryCommand(primary *Primary, success func(), fail func()) {
	if primary == nil {
		success()
	}
	table.AddCommand("DropPrimary", success, fail, primary.Name, primary.Columns)
}

// DropIndexCommand add a new command that dropping a index
func (table *Table) DropIndexCommand(name string, success func(), fail func()) {
	table.AddCommand("DropIndex", success, fail, name)
}

// RenameIndexCommand add a new command that renaming a index
func (table *Table) RenameIndexCommand(old string, new string, success func(), fail func()) {
	table.AddCommand("RenameIndex", success, fail, old, new)
}
