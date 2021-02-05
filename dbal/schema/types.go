package schema

import "github.com/jmoiron/sqlx"

// Connection DB Connection
type Connection struct{ Write *sqlx.DB }

// Schema The database Schema interface
type Schema interface {
	Create(string, func(table *Blueprint))
	Drop()
	DropIfExists()
	Rename()
	GetColumnType(string) string
	GetIndexType(string) string
}

// BlueprintAPI  the bluprint interface
type BlueprintAPI interface {
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
	BlueprintAPI
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
}

// Index  the table index definition
type Index struct {
	Comment string
	Name    string
	Type    string
	Columns []*Column
	Table   *Blueprint
}
