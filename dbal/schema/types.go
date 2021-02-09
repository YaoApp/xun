package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar"
)

// Connection DB Connection
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *dbal.Config
	Config      *dbal.DBConfig
}

// Schema The database Schema interface
type Schema interface {
	Table(string) *Table
	HasTable(string) bool
	Create(string, func(table *Table)) error
	MustCreate(string, func(table *Table)) *Table
	Drop(string) error
	MustDrop(string)
	DropIfExists(string) error
	MustDropIfExists(string)
	Rename(string, string) error
	MustRename(string, string) *Table
	Alter(string, func(table *Table)) error
	MustAlter(string, func(table *Table)) *Table
}

// Blueprint  the bluprint interface
type Blueprint interface {
	BigInteger()
	String(name string, length int) *Table
	Primary()
}

// Builder the dbal schema driver
type Builder struct {
	Conn    *Connection
	Grammar grammar.Grammar
	Schema
}

// Table the table blueprint
type Table struct {
	Blueprint
	Builder   *Builder
	Comment   string
	Name      string
	Columns   []*Column
	ColumnMap map[string]*Column
	Indexes   []*Index
	IndexMap  map[string]*Index
	alter     bool
	Table     *grammar.Table
}

// Column the table column definition
type Column struct {
	Comment  string
	Name     string
	Type     string
	Length   int
	Args     interface{}
	Default  interface{}
	Nullable *bool
	Unsigned *bool
	Table    *Table
	dropped  bool
	renamed  bool
	newname  string
}

// Index  the table index definition
type Index struct {
	Comment string
	Name    string
	Type    string
	Columns []*Column
	Table   *Table
	dropped bool
	renamed bool
	newname string
}
