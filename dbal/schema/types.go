package schema

import "github.com/jmoiron/sqlx"

// Connection DB Connection
type Connection struct{ Write *sqlx.DB }

// Schema The database Schema interface
type Schema interface {
	Table(string) *Blueprint
	HasTable(string) bool
	Create(string, func(table *Blueprint)) error
	MustCreate(string, func(table *Blueprint)) *Blueprint
	Drop(string) error
	MustDrop(string)
	DropIfExists(string) error
	MustDropIfExists(string)
	Rename(string, string) error
	MustRename(string, string) *Blueprint
	Alter(string, func(table *Blueprint)) error
	GetColumnType(string) string
	GetIndexType(string) string
}

// BlueprintInterface  the bluprint interface
type BlueprintInterface interface {
	BigInteger()
	String(name string, length int) *Blueprint
	Primary()
}

// Builder the dbal schema driver
type Builder struct {
	Conn *Connection
	Schema
}

// Blueprint the table blueprint
type Blueprint struct {
	BlueprintInterface
	Builder   *Builder
	Comment   string
	Name      string
	Columns   []*Column
	ColumnMap map[string]*Column
	Indexes   []*Index
	IndexMap  map[string]*Index
}

// Column the table column definition
type Column struct {
	Comment  string
	Name     string
	Type     string
	Length   *int
	Args     interface{}
	Default  interface{}
	Nullable *bool
	Unsigned *bool
	Table    *Blueprint
	Changed  bool
	Removed  bool
	Newname  string
}

// Index  the table index definition
type Index struct {
	Comment string
	Name    string
	Type    string
	Columns []*Column
	Table   *Blueprint
}

// TableField the table field
type TableField struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   interface{} `db:"Extra"`
}

// TableIndex the table index
type TableIndex struct {
	NonUnique  int
	KeyName    string
	SeqInIndex int
	ColumnName string
}
