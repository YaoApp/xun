package schema

import "github.com/yaoapp/xun/dbal"

// Schema The schema interface
type Schema interface {
	Get(string) (Blueprint, error)
	Create(string, func(table Blueprint)) error
	Drop(string) error
	Alter(string, func(table Blueprint)) error
	HasTable(string) bool
	Rename(string, string) error
	DropIfExists(string) error
	GetVersion() (*dbal.Version, error)

	MustGet(string) Blueprint
	MustCreate(string, func(table Blueprint)) Blueprint
	MustDrop(string)
	MustAlter(string, func(table Blueprint)) Blueprint
	MustRename(string, string) Blueprint
	MustDropIfExists(string)
	MustGetVersion() *dbal.Version
}

// Blueprint the table operating interface
type Blueprint interface {

	// defined in table.go
	GetName() string
	GetPrefix() string
	GetFullName() string

	GetColumns() map[string]*Column
	GetIndexes() map[string]*Index

	GetTable() *Table

	// defined in column.go
	NewColumn(name string) *Column
	PushColumn(column *Column) *Table
	GetColumn(name string) *Column
	Column(name string) *Column
	HasColumn(name ...string) bool
	PutColumn(column *Column) *Table
	AddColumn(column *Column) *Table
	ChangeColumn(column *Column) *Table
	RenameColumn(old string, new string) *Column
	DropColumn(name ...string)

	// defined in primry.go
	GetPrimary() *Primary
	AddPrimary(columnName ...string)
	AddPrimaryWithName(name string, columnName ...string)
	DropPrimary()

	// defined in index.go
	NewIndex(name string, columns ...*Column) *Index
	PushIndex(index *Index) *Table
	GetIndex(name string) *Index
	Index(name string) *Index
	HasIndex(name ...string) bool
	PutIndex(key string, columnNames ...string) *Table
	PutUnique(key string, columnNames ...string) *Table
	AddIndex(key string, columnNames ...string) *Table
	AddUnique(key string, columnNames ...string) *Table
	ChangeIndex(key string, columnNames ...string) *Table
	RenameIndex(old string, new string) *Index
	DropIndex(key ...string)

	// defined in constraint.go
	AddUniqueConstraint(name string, columnNames ...string)
	GetUniqueConstraint(name string)
	DropUniqueConstraint(name string)

	// defined in blueprint.go
	// Character types
	String(name string, args ...int) *Column
	Char(name string, args ...int) *Column
	Text(name string) *Column
	MediumText(name string) *Column
	LongText(name string) *Column

	// Binary types
	Binary(name string, args ...int) *Column

	// Date time types
	Date(name string) *Column
	DateTime(name string, args ...int) *Column
	DateTimeTz(name string, args ...int) *Column
	Time(name string, args ...int) *Column
	TimeTz(name string, args ...int) *Column
	Timestamp(name string, args ...int) *Column
	TimestampTz(name string, args ...int) *Column

	// Numberic types
	TinyInteger(name string) *Column
	UnsignedTinyInteger(name string) *Column
	TinyIncrements(name string) *Column

	SmallInteger(name string) *Column
	UnsignedSmallInteger(name string) *Column
	SmallIncrements(name string) *Column

	Integer(name string) *Column
	UnsignedInteger(name string) *Column
	Increments(name string) *Column

	BigInteger(name string) *Column
	UnsignedBigInteger(name string) *Column
	BigIncrements(name string) *Column
	ID(name string) *Column
	ForeignID(name string) *Column

	Decimal(name string, args ...int) *Column
	UnsignedDecimal(name string, args ...int) *Column

	Float(name string, args ...int) *Column
	UnsignedFloat(name string, args ...int) *Column

	Double(name string, args ...int) *Column
	UnsignedDouble(name string, args ...int) *Column

	// boolean, enum types
	Boolean(name string) *Column
	Enum(name string, option []string) *Column

	// json, jsonb types
	JSON(name string) *Column
	JSONB(name string) *Column

	// uuid, geometry, geometryCollection, point, multiPoint, polygon, multiPolygon
	UUID(name string) *Column

	// ipAddress, macAddress, year

	// timestamps, timestampsTz, nullableTimestamps,  softDeletes, softDeletesTz

	// morphs, nullableMorphs, uuidMorphs nullableUuidMorphs

}
