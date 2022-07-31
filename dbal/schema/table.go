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
		Name:        name,
		Prefix:      builder.Conn.Option.Prefix,
		Table:       dbal.NewTable(tableName, builder.Schema, builder.Database),
		Builder:     builder,
		IndexNames:  []string{},
		ColumnNames: []string{},
		ColumnMap:   map[string]*Column{},
		IndexMap:    map[string]*Index{},
	}
	return table
}

// GetColumnNames Get the column names
func (table *Table) GetColumnNames() []string {
	return table.ColumnNames
}

// GetColumns Get the columns map of the table
func (table *Table) GetColumns() map[string]*Column {
	return table.ColumnMap
}

// GetIndexNames Get the index names
func (table *Table) GetIndexNames() []string {
	return table.IndexNames
}

// GetIndexes Get the indexes map of the table
func (table *Table) GetIndexes() map[string]*Index {
	return table.IndexMap
}

// Get Get the DBAL table instance
func (table *Table) Get() *Table {
	return table
}
