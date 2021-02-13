package schema

import "fmt"

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
	index := column.Table.NewIndex(name, column)
	index.Type = "unique"
	if column.Table.HasIndex(name) {
		column.Table.DropIndexCommand(name)
	}
	column.Table.CreateIndexCommand(&index.Index)
	return column
}

// Primary set as primary key
func (column *Column) Primary() *Column {
	name := fmt.Sprintf("%s_%s", column.Name, "primary")
	if column.HasIndex(name) {
		return column
	}
	index := column.Table.NewIndex(name, column)
	index.Type = "primary"
	column.Column.Primary = true
	if column.Table.HasIndex(name) {
		column.Table.DropIndexCommand(name)
	}
	column.Table.CreateIndexCommand(&index.Index)
	column.Table.Primary = &column.Column
	column.NotNull()
	return column
}

// Index set as index key
func (column *Column) Index() *Column {
	name := fmt.Sprintf("%s_%s", column.Name, "index")
	if column.HasIndex(name) {
		return column
	}
	index := column.Table.NewIndex(name, column)
	index.Type = "index"
	if column.Table.HasIndex(name) {
		column.Table.DropIndexCommand(name)
	}
	column.Table.CreateIndexCommand(&index.Index)
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
	column.Extra = "AutoIncrement"
	return column
}
