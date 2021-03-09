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
// @Todo:
//   1. tinyInteger() Create a new tiny integer (1-byte) column on the table.
//   2. MediumInteger()  Create a new medium integer (3-byte) column on the table.

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

// Integer Create a new integer (4-byte) column on the table.
func (table *Table) Integer(name string) *Column {
	column := table.NewColumn(name).SetType("integer")
	table.PutColumn(column)
	return column
}

// UnsignedInteger Create a new auto-incrementing integer (4-byte) column on the table.
func (table *Table) UnsignedInteger(name string) *Column {
	return table.Integer(name).Unsigned()
}

// Increments Create a new auto-incrementing big integer (2-byte) column on the table.
func (table *Table) Increments(name string) *Column {
	return table.UnsignedInteger(name).AutoIncrement()
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

// Decimal Create a new decimal (16-byte) column on the table.
func (table *Table) Decimal(name string, total int, places int) *Column {
	column := table.NewColumn(name).
		SetType("decimal").
		SetPrecision(total).
		SetScale(places)
	table.PutColumn(column)
	return column
}

// UnsignedDecimal Create a new unsigned decimal (16-byte) column on the table.
func (table *Table) UnsignedDecimal(name string, total int, places int) *Column {
	return table.Decimal(name, total, places).Unsigned()
}

// Float Create a new float (4-byte) column on the table.
func (table *Table) Float(name string, total int, places int) *Column {
	column := table.NewColumn(name).
		SetType("float").
		SetPrecision(total).
		SetScale(places)
	table.PutColumn(column)
	return column
}

// UnsignedFloat Create a new unsigned float (4-byte) column on the table.
func (table *Table) UnsignedFloat(name string, total int, places int) *Column {
	return table.Float(name, total, places).Unsigned()
}
