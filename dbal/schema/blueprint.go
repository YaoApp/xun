package schema

// GetName get the table name
func (table *Table) GetName() string {
	return table.Name
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

// GetColumns Get the columns map of the table
func (table *Table) GetColumns() map[string]*Column {
	return table.ColumnMap
}

// GetIndexes Get the indexes map of the table
func (table *Table) GetIndexes() map[string]*Index {
	return table.IndexMap
}

// GetIndex get the index instance for the given name, create if not exists.
func (table *Table) GetIndex(name string) *Index {
	index, has := table.IndexMap[name]
	if !has {
		index = table.NewIndex(name)
	}
	return index
}

// GetColumn get the column instance for the given name, create if not exists.
func (table *Table) GetColumn(name string) *Column {
	column, has := table.ColumnMap[name]
	if !has {
		column = table.NewColumn(name)
	}
	return column
}

// AddColumn add or modify a column to the table
func (table *Table) AddColumn(column *Column) *Column {
	if table.HasColumn(column.Name) {
		table.ModifyColumnCommand(column.Column)
		table.onChange("ModifyColumn", column)
		return column
	}
	table.AddColumnCommand(column.Column)
	table.onChange("AddColumn", column)
	return column
}

// DropColumn Indicate that the given columns should be dropped.
func (table *Table) DropColumn(name ...string) {
	for _, n := range name {
		table.DropColumnCommand(n)
	}
	table.onChange("DropColumn", name)
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Table) RenameColumn(old string, new string) *Column {
	table.RenameColumnCommand(old, new)
	column := table.GetColumn(old)
	column.Name = new
	table.onChange("RenameColumn", old, new)
	return column
}

// String Create a new string column on the table.
func (table *Table) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = &length
	column.Type = "string"
	table.AddColumn(column)
	return column
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) BigInteger(name string) *Column {
	column := table.NewColumn(name)
	column.Type = "bigInteger"
	table.AddColumn(column)
	return column
}

// UnsignedBigInteger Create a new unsigned big integer (8-byte) column on the table.
func (table *Table) UnsignedBigInteger(name string) *Column {
	return table.BigInteger(name).Unsigned()
}

// BigIncrements Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) BigIncrements(name string) *Column {
	return table.UnsignedBigInteger(name).AutoIncrement()
}

// ID Alias BigIncrements. Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) ID(name string) *Column {
	return table.BigIncrements(name).Primary()
}
