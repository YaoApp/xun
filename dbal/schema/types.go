package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar"
)

// Schema The schema interface
type Schema interface {
	SetMode(string) Schema
	Table(string) *Table
	HasTable(string) bool
	Create(string, func(table Blueprint)) error
	MustCreate(string, func(table Blueprint)) *Table
	Drop(string) error
	MustDrop(string)
	DropIfExists(string) error
	MustDropIfExists(string)
	Rename(string, string) error
	MustRename(string, string) *Table
	Alter(string, func(table Blueprint)) error
	MustAlter(string, func(table Blueprint)) *Table
}

// Blueprint the table operating interface
type Blueprint interface {
	String(name string, length int) *Column
	HasColumn(name ...string) bool
	DropColumn(name ...string)
	RenameColumn(old string, new string) *Column
	HasIndex(name ...string) bool
	DropIndex(name ...string)
	RenameIndex(old string, new string) *Index
}

// Connection the database connection for schema operating
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *dbal.Config
	Config      *dbal.DBConfig
}

// Builder the table schema builder struct
type Builder struct {
	Conn    *Connection
	Grammar grammar.Grammar
	Mode    string
}

// Table the table struct
type Table struct {
	grammar.Table
	Builder   *Builder
	ColumnMap map[string]*Column
	IndexMap  map[string]*Index
}

// Column the table column struct
type Column struct {
	grammar.Column
	Table *Table
}

// Index  the table index struct
type Index struct {
	grammar.Index
	Table *Table
}
