package schema

import (
	"github.com/yaoapp/xun/dbal"
)

// GetIndex get the index instance for the given name,if the index does not exist return nil.
func (table *Table) GetIndex(name string) *Index {
	return table.IndexMap[name]
}

// HasIndex Determine if the table has a given index.
func (table *Table) HasIndex(name ...string) bool {
	has := true
	for _, n := range name {
		_, has = table.IndexMap[n]
		if !has {
			return has
		}
	}
	return has
}

// AddIndex Indicate that the given index should be created.
func (table *Table) AddIndex(key string, columnNames ...string) *Table {
	columns := []*Column{}
	for _, name := range columnNames {
		columns = append(columns, table.GetColumn(name))
	}
	index := table.newIndex(key, columns...)
	index.Type = "index"
	table.pushIndex(index)
	table.createIndexCommand(index.Index, nil, func() {
		delete(table.IndexMap, index.Name)
	})
	return table
}

// AddUnique Indicate that the given unique index should be created.
func (table *Table) AddUnique(key string, columnNames ...string) *Table {
	columns := []*Column{}
	for _, name := range columnNames {
		columns = append(columns, table.GetColumn(name))
	}
	index := table.newIndex(key, columns...)
	index.Type = "unique"
	table.pushIndex(index)
	table.createIndexCommand(index.Index, nil, func() {
		delete(table.IndexMap, index.Name)
	})
	return table
}

// DropIndex Indicate that the given indexes should be dropped.
func (table *Table) DropIndex(key ...string) {
	for _, n := range key {
		table.dropIndexCommand(n, func() {
			delete(table.IndexMap, n)
		}, nil)
	}
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Table) RenameIndex(old string, new string) *Index {
	index := table.GetIndex(old)
	index.Name = new
	table.IndexMap[new] = index
	table.renameIndexCommand(old, new,
		func() {
			delete(table.IndexMap, old)
		},
		func() {
			delete(table.IndexMap, new)
		})
	return index
}

// newIndex Create a new index instance
func (table *Table) newIndex(name string, columns ...*Column) *Index {
	cols := []*dbal.Column{}
	for _, column := range columns {
		cols = append(cols, column.Column)
	}
	index := &Index{
		Index: table.Table.NewIndex(name, cols...),
		Table: table,
	}

	// mapping index
	for _, column := range columns {
		column.Indexes = append(column.Indexes, index.Index)
	}
	return index
}

// pushIndex add an index to the table
func (table *Table) pushIndex(index *Index) *Table {
	table.Table.PushIndex(index.Index)
	table.IndexMap[index.Name] = index
	return table
}
