package schema

// Character types

// String Create a new string column on the table.
func (table *Table) String(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("string")
	column.MaxLength = 65535
	column.DefaultLength = 200
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.PutColumn(column)
	return column
}

// Char Create a new char column on the table.
func (table *Table) Char(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("char")
	column.MaxLength = 30
	column.DefaultLength = 10
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.PutColumn(column)
	return column
}

// Text Create a new text column on the table.
func (table *Table) Text(name string) *Column {
	column := table.NewColumn(name).SetType("text")
	table.PutColumn(column)
	return column
}

// MediumText Create a new medium text column on the table.
func (table *Table) MediumText(name string) *Column {
	column := table.NewColumn(name).SetType("mediumText")
	table.PutColumn(column)
	return column
}

// LongText Create a new long text column on the table.
func (table *Table) LongText(name string) *Column {
	column := table.NewColumn(name).SetType("longText")
	table.PutColumn(column)
	return column
}

// Binary types

// Binary Create a new binary column on the table.
func (table *Table) Binary(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("binary")
	column.MaxLength = 65535
	column.DefaultLength = 255
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.PutColumn(column)
	return column
}

// Date time types

// Date Create a new date column on the table.
func (table *Table) Date(name string) *Column {
	column := table.NewColumn(name).SetType("date")
	table.PutColumn(column)
	return column
}

// DateTime Create a new date-time column on the table.
func (table *Table) DateTime(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("dateTime")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// DateTimeTz Create a new date-time column (with time zone) on the table.
func (table *Table) DateTimeTz(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("dateTimeTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// Time Create a new time column on the table.
func (table *Table) Time(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("time")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// TimeTz Create a new time column (with time zone) on the table.
func (table *Table) TimeTz(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("timeTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// Timestamp Create a new timestamp column on the table.
func (table *Table) Timestamp(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("timestamp")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// TimestampTz Create a new timestamp (with time zone) column on the table.
func (table *Table) TimestampTz(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("timestampTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.PutColumn(column)
	return column
}

// Numberic types
// @Todo:
//   1. tinyInteger() Create a new tiny integer (1-byte) column on the table. [done]
//   2. MediumInteger()  Create a new medium integer (3-byte) column on the table.

// TinyInteger Create a new tiny integer (1-byte) column on the table.
func (table *Table) TinyInteger(name string) *Column {
	column := table.NewColumn(name).SetType("tinyInteger")
	table.PutColumn(column)
	return column
}

// UnsignedTinyInteger Create a new auto-incrementing tiny integer (1-byte) column on the table.
func (table *Table) UnsignedTinyInteger(name string) *Column {
	return table.TinyInteger(name).Unsigned()
}

// TinyIncrements Create a new auto-incrementing tiny integer (1-byte) column on the table.
func (table *Table) TinyIncrements(name string) *Column {
	return table.UnsignedTinyInteger(name).AutoIncrement()
}

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

// SmallIncrements Create a new auto-incrementing small integer (2-byte) column on the table.
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

// ForeignID Alias UnsignedBigInteger. Create a new unsigned big integer (8-byte) column on the table.
func (table *Table) ForeignID(name string) *Column {
	return table.UnsignedBigInteger(name)
}

// Decimal Create a new decimal (16-byte) column on the table.
func (table *Table) Decimal(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("decimal")
	column.MaxPrecision = 65
	column.MaxScale = 30
	column.DefaultPrecision = 10
	column.DefaultScale = 2

	total := column.DefaultPrecision
	places := column.DefaultScale
	if len(args) >= 1 {
		total = args[0]
	}
	if len(args) >= 2 {
		places = args[1]
	}
	column.SetPrecision(total).SetScale(places)
	table.PutColumn(column)
	return column
}

// UnsignedDecimal Create a new unsigned decimal (16-byte) column on the table.
func (table *Table) UnsignedDecimal(name string, args ...int) *Column {
	return table.Decimal(name, args...).Unsigned()
}

// Float Create a new float (4-byte) column on the table.
func (table *Table) Float(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("float")
	column.MaxPrecision = 23
	column.DefaultPrecision = 10
	column.MaxScale = 22
	column.DefaultScale = 2

	total := column.DefaultPrecision
	places := column.DefaultScale
	if len(args) >= 1 {
		total = args[0]
	}
	if len(args) >= 2 {
		places = args[1]
	}
	column.SetPrecision(total).SetScale(places)
	table.PutColumn(column)
	return column
}

// UnsignedFloat Create a new unsigned float (4-byte) column on the table.
func (table *Table) UnsignedFloat(name string, args ...int) *Column {
	return table.Float(name, args...).Unsigned()
}

// Double Create a new double (8-byte) column on the table.
func (table *Table) Double(name string, args ...int) *Column {
	column := table.NewColumn(name).SetType("double")
	column.MaxPrecision = 53
	column.MaxScale = 52
	column.DefaultPrecision = 24
	column.DefaultScale = 2

	total := column.DefaultPrecision
	places := column.DefaultScale
	if len(args) >= 1 {
		total = args[0]
	}
	if len(args) >= 2 {
		places = args[1]
	}
	column.SetPrecision(total).SetScale(places)
	table.PutColumn(column)
	return column
}

// UnsignedDouble Create a new unsigned double (8-byte) column on the table.
func (table *Table) UnsignedDouble(name string, args ...int) *Column {
	return table.Double(name, args...).Unsigned()
}

// Boolean Create a new boolean column on the table.
func (table *Table) Boolean(name string) *Column {
	column := table.NewColumn(name).SetType("boolean")
	table.PutColumn(column)
	return column
}

// Enum Create a new enum column on the table.
func (table *Table) Enum(name string, option []string) *Column {
	column := table.NewColumn(name).SetType("enum")
	column.Option = option
	table.PutColumn(column)
	return column
}

// JSON Create a new json column on the table.
func (table *Table) JSON(name string) *Column {
	column := table.NewColumn(name).SetType("json")
	table.PutColumn(column)
	return column
}

// JSONB  Create a new jsonb column on the table.
func (table *Table) JSONB(name string) *Column {
	column := table.NewColumn(name).SetType("jsonb")
	table.PutColumn(column)
	return column
}

// UUID Create a new uuid column on the table.
func (table *Table) UUID(name string) *Column {
	column := table.NewColumn(name).SetType("uuid")
	table.PutColumn(column)
	return column
}

// IPAddress Create a new IP address ( integer 4-byte ) column on the table.
func (table *Table) IPAddress(name string) *Column {
	column := table.NewColumn(name).SetType("ipAddress")
	table.PutColumn(column)
	return column
}
