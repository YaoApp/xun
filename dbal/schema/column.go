package schema

import (
	"fmt"

	"github.com/yaoapp/xun/utils"
)

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
	column.Table.CreateUnique(name, column.Name)
	return column
}

// Primary set as primary key
func (column *Column) Primary() *Column {
	if column.Column.Primary {
		return column
	}

	column.Column.Primary = true
	column.Table.CreatePrimary(column.Name)
	return column
}

// Index set as index key
func (column *Column) Index() *Column {
	name := fmt.Sprintf("%s_%s", column.Name, "index")
	if column.HasIndex(name) {
		return column
	}
	column.Table.CreateIndex(name, column.Name)
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
