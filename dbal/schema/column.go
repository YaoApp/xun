package schema

import (
	"fmt"

	"github.com/yaoapp/xun/utils"
)

// GetColumn get the column instance of the table, if the column does not exist return nil.
func (table *Table) GetColumn(name string) *Column {
	return table.ColumnMap[name]
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

// DropColumn Indicate that the given columns should be dropped.
func (table *Table) DropColumn(name ...string) {
	for _, n := range name {
		table.dropColumnCommand(n, func() {
			delete(table.ColumnMap, n)
		}, nil)
	}
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Table) RenameColumn(old string, new string) *Column {
	column := table.GetColumn(old)
	column.Name = new
	table.ColumnMap[new] = column
	table.renameColumnCommand(old, new, func() {
		delete(table.ColumnMap, old)
	}, func() {
		delete(table.ColumnMap, new)
	})
	return column
}

// newColumn Create a new column instance
func (table *Table) newColumn(name string) *Column {
	column := &Column{
		Column: table.Table.NewColumn(name),
		Table:  table,
	}
	return column
}

// pushColumn add a column to the table
func (table *Table) pushColumn(column *Column) *Table {
	table.Table.PushColumn(column.Column)
	table.ColumnMap[column.Name] = column
	return table
}

// putColumn add or modify a column to the table
func (table *Table) putColumn(column *Column) *Table {
	if table.HasColumn(column.Name) {
		table.changeColumn(column)
	} else {
		table.addColumn(column)
	}
	return table
}

// addColumn add a column to the table
func (table *Table) addColumn(column *Column) *Table {
	table.pushColumn(column)
	table.addColumnCommand(column.Column, nil, func() {
		delete(table.ColumnMap, column.Name)
	})
	return table
}

// changeColumn modify a column to the table
func (table *Table) changeColumn(column *Column) *Table {
	table.changeColumnCommand(column.Column, func() {
		table.ColumnMap[column.Name] = column
	}, nil)
	return table
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

// the column portables

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

// SetLength set the column Length attribute to the given length
func (column *Column) SetLength(length int) *Column {
	if column.MaxLength == 0 {
		return column
	}
	if length > column.MaxLength || length == 0 {
		length = column.DefaultLength
	}
	column.Length = &length
	return column
}

// SetType set the column type attribute to the given type name
func (column *Column) SetType(typ string) *Column {
	column.Type = typ
	return column
}

// SetComment set the column comment to the given value
func (column *Column) SetComment(comment string) *Column {
	column.Comment = &comment
	return column
}

// SetDefault set the column default attribute to the given type name
func (column *Column) SetDefault(v interface{}) *Column {
	column.Default = v
	return column
}

// SetDateTimePrecision set the column precision to the given value
func (column *Column) SetDateTimePrecision(precision int) *Column {
	if column.MaxDateTimePrecision == 0 {
		return column
	}
	if precision > column.MaxDateTimePrecision || precision == 0 {
		precision = column.DefaultDateTimePrecision
	}
	column.DateTimePrecision = &precision
	return column
}

// SetPrecision set the column precision to the given value
func (column *Column) SetPrecision(precision int) *Column {
	if column.MaxPrecision == 0 {
		return column
	}
	if precision > column.MaxPrecision || precision == 0 {
		precision = column.DefaultPrecision
	}
	if column.Scale != nil && *column.Scale+precision > column.MaxPrecision {
		precision = column.MaxPrecision - *column.Scale
	}

	column.Precision = &precision
	return column
}

// SetScale set the column scale to the given value
func (column *Column) SetScale(scale int) *Column {
	if column.DefaultScale == 0 {
		return column
	}
	if scale > column.MaxScale || scale == 0 {
		scale = column.DefaultScale
	}
	if column.Precision != nil && *column.Precision+scale > column.MaxPrecision {
		scale = column.MaxPrecision - *column.Precision
	}

	if column.Precision != nil && scale > *column.Precision {
		scale = *column.Precision
	}
	column.Scale = &scale
	return column
}
