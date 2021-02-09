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

// Builder the dbal schema driver
type Builder struct {
	Conn    *Connection
	Grammar grammar.Grammar
	Schema
}

// Table the table blueprint
type Table struct {
	grammar.Table
	Builder   *Builder
	ColumnMap map[string]*Column
	IndexMap  map[string]*Index
}

// Column the table column definition
type Column struct {
	grammar.Column
	Table *Table
}

// Index  the table index definition
type Index struct {
	grammar.Index
	Table *Table
}
