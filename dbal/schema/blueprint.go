package schema

// String Create a new string column on the table.
func (table *Table) String(name string, length int) *Column {
	column := table.NewColumn(name).
		SetLength(length).
		SetType("string")

	table.PutColumn(column)
	return column
}

// Numberic types

// SmallInteger Create a new small integer (2-byte) column on the table.
func (table *Table) SmallInteger(name string) *Column {
	column := table.NewColumn(name).SetType("smallInteger")
	table.PutColumn(column)
	return column
}

// UnsignedSmallInteger Create a new unsigned small integer (2-byte) column on the table.
func (table *Table) UnsignedSmallInteger(name string) *Column {
	return table.SmallInteger(name).Unsigned()
}

// SmallIncrements Create a new auto-incrementing big integer (2-byte) column on the table.
func (table *Table) SmallIncrements(name string) *Column {
	return table.UnsignedSmallInteger(name).AutoIncrement()
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) BigInteger(name string) *Column {
	column := table.NewColumn(name).SetType("bigInteger")
	table.PutColumn(column)
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
