package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
)

// Schema The schema interface
type Schema interface {
	SetOption(option *dbal.Option)

	Builder() *Builder
	GetConnection() (*dbal.Connection, error)
	GetDB() (*sqlx.DB, error)
	GetVersion() (*dbal.Version, error)

	GetTables() ([]string, error)

	GetTable(name string) (Blueprint, error)
	CreateTable(name string, createFunc func(table Blueprint)) error
	DropTable(name string) error
	AlterTable(name string, alterFunc func(table Blueprint)) error
	HasTable(name string) (bool, error)
	RenameTable(old string, new string) error
	DropTableIfExists(name string) error

	MustGetConnection() *dbal.Connection
	MustGetDB() *sqlx.DB
	MustGetVersion() *dbal.Version

	MustGetTables() []string

	MustGetTable(name string) Blueprint
	MustCreateTable(name string, createFunc func(table Blueprint))
	MustDropTable(name string)
	MustAlterTable(name string, alterFunc func(table Blueprint))
	MustHasTable(name string) bool
	MustRenameTable(old string, new string) Blueprint
	MustDropTableIfExists(name string)

	DB() *sqlx.DB // alias MustGetDB
}

// Blueprint the table operating interface
type Blueprint interface {

	// defined in table.go
	Get() *Table
	GetName() string
	GetPrefix() string
	GetFullName() string
	GetColumnNames() []string
	GetColumns() map[string]*Column
	GetIndexNames() []string
	GetIndexes() map[string]*Index

	// defined in column.go
	GetColumn(name string) *Column
	HasColumn(name ...string) bool
	RenameColumn(old string, new string) *Column
	DropColumn(name ...string)

	// defined in primry.go
	GetPrimary() *Primary
	AddPrimary(columnName ...string)
	DropPrimary()

	// defined in index.go
	GetIndex(name string) *Index
	HasIndex(name ...string) bool
	AddIndex(name string, columnNames ...string) *Table
	AddUnique(name string, columnNames ...string) *Table
	AddFulltext(name string, columnNames ...string) *Table
	RenameIndex(old string, new string) *Index
	DropIndex(name ...string)

	// defined in constraint.go
	// @todo: GetUniqueConstraint, AddUniqueConstraint, DropUniqueConstraint

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
	// @todo: MediumInteger
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

	// uuid, ipAddress, macAddress, year etc.
	// @todo: geometry, geometryCollection, point, multiPoint, polygon, multiPolygon
	UUID(name string) *Column
	IPAddress(name string) *Column
	MACAddress(name string) *Column
	Year(name string) *Column

	// timestamps, timestampsTz,DropTimestamps, DropTimestampsTz, softDeletes, softDeletesTz, DropSoftDeletes, DropSoftDeletesTz
	Timestamps(args ...int) map[string]*Column
	TimestampsTz(args ...int) map[string]*Column
	DropTimestamps()
	DropTimestampsTz()

	SoftDeletes(args ...int) *Column
	SoftDeletesTz(args ...int) *Column
	DropSoftDeletes()
	DropSoftDeletesTz()

	//@todo: morphs, nullableMorphs, uuidMorphs nullableUuidMorphs

}
