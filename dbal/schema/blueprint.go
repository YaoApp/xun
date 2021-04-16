package schema

// Character types

// String Create a new string column on the table.
func (table *Table) String(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("string")
	column.MaxLength = 65535
	column.DefaultLength = 200
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.putColumn(column)
	return column
}

// Char Create a new char column on the table.
func (table *Table) Char(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("char")
	column.MaxLength = 30
	column.DefaultLength = 10
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.putColumn(column)
	return column
}

// Text Create a new text column on the table.
func (table *Table) Text(name string) *Column {
	column := table.newColumn(name).SetType("text")
	table.putColumn(column)
	return column
}

// MediumText Create a new medium text column on the table.
func (table *Table) MediumText(name string) *Column {
	column := table.newColumn(name).SetType("mediumText")
	table.putColumn(column)
	return column
}

// LongText Create a new long text column on the table.
func (table *Table) LongText(name string) *Column {
	column := table.newColumn(name).SetType("longText")
	table.putColumn(column)
	return column
}

// Binary types

// Binary Create a new binary column on the table.
func (table *Table) Binary(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("binary")
	column.MaxLength = 65535
	column.DefaultLength = 255
	length := column.DefaultLength
	if len(args) >= 1 {
		length = args[0]
	}
	column.SetLength(length)
	table.putColumn(column)
	return column
}

// Date time types

// Date Create a new date column on the table.
func (table *Table) Date(name string) *Column {
	column := table.newColumn(name).SetType("date")
	table.putColumn(column)
	return column
}

// DateTime Create a new date-time column on the table.
func (table *Table) DateTime(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("dateTime")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// DateTimeTz Create a new date-time column (with time zone) on the table.
func (table *Table) DateTimeTz(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("dateTimeTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// Time Create a new time column on the table.
func (table *Table) Time(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("time")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// TimeTz Create a new time column (with time zone) on the table.
func (table *Table) TimeTz(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("timeTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// Timestamp Create a new timestamp column on the table.
func (table *Table) Timestamp(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("timestamp")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// TimestampTz Create a new timestamp (with time zone) column on the table.
func (table *Table) TimestampTz(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("timestampTz")
	column.MaxDateTimePrecision = 6
	column.DefaultDateTimePrecision = 0
	precision := column.DefaultDateTimePrecision
	if len(args) >= 1 {
		precision = args[0]
	}
	column.SetDateTimePrecision(precision)
	table.putColumn(column)
	return column
}

// Numberic types
// @Todo:
//   1. tinyInteger() Create a new tiny integer (1-byte) column on the table. [done]
//   2. MediumInteger()  Create a new medium integer (3-byte) column on the table.

// TinyInteger Create a new tiny integer (1-byte) column on the table.
func (table *Table) TinyInteger(name string) *Column {
	column := table.newColumn(name).SetType("tinyInteger")
	table.putColumn(column)
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
	column := table.newColumn(name).SetType("smallInteger")
	table.putColumn(column)
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
	column := table.newColumn(name).SetType("integer")
	table.putColumn(column)
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
	column := table.newColumn(name).SetType("bigInteger")
	table.putColumn(column)
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
	column := table.newColumn(name).SetType("decimal")
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
	table.putColumn(column)
	return column
}

// UnsignedDecimal Create a new unsigned decimal (16-byte) column on the table.
func (table *Table) UnsignedDecimal(name string, args ...int) *Column {
	return table.Decimal(name, args...).Unsigned()
}

// Float Create a new float (4-byte) column on the table.
func (table *Table) Float(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("float")
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
	table.putColumn(column)
	return column
}

// UnsignedFloat Create a new unsigned float (4-byte) column on the table.
func (table *Table) UnsignedFloat(name string, args ...int) *Column {
	return table.Float(name, args...).Unsigned()
}

// Double Create a new double (8-byte) column on the table.
func (table *Table) Double(name string, args ...int) *Column {
	column := table.newColumn(name).SetType("double")
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
	table.putColumn(column)
	return column
}

// UnsignedDouble Create a new unsigned double (8-byte) column on the table.
func (table *Table) UnsignedDouble(name string, args ...int) *Column {
	return table.Double(name, args...).Unsigned()
}

// Boolean Create a new boolean column on the table.
func (table *Table) Boolean(name string) *Column {
	column := table.newColumn(name).SetType("boolean")
	table.putColumn(column)
	return column
}

// Enum Create a new enum column on the table.
func (table *Table) Enum(name string, option []string) *Column {
	column := table.newColumn(name).SetType("enum")
	column.Option = option
	table.putColumn(column)
	return column
}

// JSON Create a new json column on the table.
func (table *Table) JSON(name string) *Column {
	column := table.newColumn(name).SetType("json")
	table.putColumn(column)
	return column
}

// JSONB  Create a new jsonb column on the table.
func (table *Table) JSONB(name string) *Column {
	column := table.newColumn(name).SetType("jsonb")
	table.putColumn(column)
	return column
}

// UUID Create a new uuid column on the table.
func (table *Table) UUID(name string) *Column {
	column := table.newColumn(name).SetType("uuid")
	table.putColumn(column)
	return column
}

// IPAddress Create a new IP address ( integer 4-byte ) column on the table.
func (table *Table) IPAddress(name string) *Column {
	column := table.newColumn(name).SetType("ipAddress")
	table.putColumn(column)
	return column
}

// MACAddress Create a new MAC address column on the table.
func (table *Table) MACAddress(name string) *Column {
	column := table.newColumn(name).SetType("macAddress")
	table.putColumn(column)
	return column
}

// Year Create a new year column on the table.
func (table *Table) Year(name string) *Column {
	column := table.newColumn(name).SetType("year")
	table.putColumn(column)
	return column
}

// Timestamps Add nullable creation and update timestamps to the table.
func (table *Table) Timestamps(args ...int) map[string]*Column {
	return map[string]*Column{
		"created_at": table.Timestamp("created_at", args...).NotNull().SetDefaultRaw("NOW()").Index(),
		"updated_at": table.Timestamp("updated_at", args...).Null().Index(),
	}
}

// TimestampsTz Add creation and update timestampTz columns to the table.
func (table *Table) TimestampsTz(args ...int) map[string]*Column {
	return map[string]*Column{
		"created_at": table.TimestampTz("created_at", args...).NotNull().SetDefaultRaw("NOW()").Index(),
		"updated_at": table.TimestampTz("updated_at", args...).Null().Index(),
	}
}

// DropTimestamps drop the "created_at", "updated_at" timestamp columns.
func (table *Table) DropTimestamps() {
	table.DropColumn("created_at", "updated_at")
}

// DropTimestampsTz drop the "created_at", "updated_at" timestamp columns.
func (table *Table) DropTimestampsTz() {
	table.DropTimestamps()
}

// SoftDeletes Add a "deleted_at" timestamp for the table.
func (table *Table) SoftDeletes(args ...int) *Column {
	return table.Timestamp("deleted_at", args...).Null().Index()
}

// SoftDeletesTz Add a "deleted_at" timestampTz for the table.
func (table *Table) SoftDeletesTz(args ...int) *Column {
	return table.TimestampTz("deleted_at", args...).Null().Index()
}

// DropSoftDeletes drop the "deleted_at" timestamp columns.
func (table *Table) DropSoftDeletes() {
	table.DropColumn("deleted_at")
}

// DropSoftDeletesTz drop the "deleted_at" timestamp columns.
func (table *Table) DropSoftDeletesTz() {
	table.DropSoftDeletes()
}
