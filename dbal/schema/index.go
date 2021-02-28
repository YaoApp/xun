package schema

import (
	"github.com/yaoapp/xun/dbal"
)

// NewIndex Create a new index instance
func (table *Table) NewIndex(name string, columns ...*Column) *Index {
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

// PushIndex add an index to the table
func (table *Table) PushIndex(index *Index) *Table {
	table.Table.PushIndex(index.Index)
	table.IndexMap[index.Name] = index
	return table
}

// GetIndex get the index instance for the given name,if the index does not exist return nil.
func (table *Table) GetIndex(name string) *Index {
	return table.IndexMap[name]
}

// Index get the index instance for the given name,if the index does not exist create a new one
func (table *Table) Index(name string) *Index {
	index, has := table.IndexMap[name]
	if !has {
		index = table.NewIndex(name)
	}
	return index
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

// PutIndex add or modify a index of the table
func (table *Table) PutIndex(key string, columnNames ...string) *Table {
	if table.HasIndex(key) {
		table.ChangeIndex(key, columnNames...)
	} else {
		table.AddIndex(key, columnNames...)
	}
	return table
}

// PutUnique add or modify a unique index of the table
func (table *Table) PutUnique(key string, columnNames ...string) *Table {
	if table.HasIndex(key) {
		table.ChangeUnique(key, columnNames...)
	} else {
		table.AddUnique(key, columnNames...)
	}
	return table
}

// AddIndex Indicate that the given index should be created.
func (table *Table) AddIndex(key string, columnNames ...string) *Table {
	columns := []*Column{}
	for _, name := range columnNames {
		columns = append(columns, table.GetColumn(name))
	}
	index := table.NewIndex(key, columns...)
	index.Type = "index"
	table.PushIndex(index)
	table.CreateIndexCommand(index.Index, nil, func() {
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
	index := table.NewIndex(key, columns...)
	index.Type = "unique"
	table.PushIndex(index)
	table.CreateIndexCommand(index.Index, nil, func() {
		delete(table.IndexMap, index.Name)
	})
	return table
}

// ChangeIndex Indicate that the given index should be changed.
func (table *Table) ChangeIndex(key string, columnNames ...string) *Table {
	return table
}

// ChangeUnique Indicate that the given unique index should be changed.
func (table *Table) ChangeUnique(key string, columnNames ...string) *Table {
	return table
}

// DropIndex Indicate that the given indexes should be dropped.
func (table *Table) DropIndex(key ...string) {
	for _, n := range key {
		table.DropIndexCommand(n, func() {
			delete(table.IndexMap, n)
		}, nil)
	}
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Table) RenameIndex(old string, new string) *Index {
	index := table.GetIndex(old)
	index.Name = new
	table.IndexMap[new] = index
	table.RenameIndexCommand(old, new,
		func() {
			delete(table.IndexMap, old)
		},
		func() {
			delete(table.IndexMap, new)
		})
	return index
}
