package schema

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
		column := table.GetColumn(n)
		column.Drop()
	}
	table.onChange("DropColumn", name)
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Table) RenameColumn(old string, new string) *Column {
	column := table.GetColumn(old)
	column.Rename(new)
	table.onChange("RenameColumn", old, new)
	return column
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

// DropIndex Indicate that the given indexes should be dropped.
func (table *Table) DropIndex(name ...string) {
	for _, n := range name {
		index := table.GetIndex(n)
		index.Drop()
	}
	table.onChange("DropIndex", name)
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Table) RenameIndex(old string, new string) *Index {
	index := table.GetIndex(old)
	index.Rename(new)
	table.onChange("RenameIndex", old, new)
	return index
}

// GetColumnListing Get the column listing for the table
func (table *Table) GetColumnListing() {
}

// GetIndexListing Get the index listing for the table
func (table *Table) GetIndexListing() {
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) BigInteger() {}

// String Create a new string column on the table.
func (table *Table) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = length
	column.Type = "string"
	table.AddColumn(column)
	return column
}
