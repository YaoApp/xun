package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar"
)

// Schema The schema interface
type Schema interface {
	Get(string) (Blueprint, error)
	Create(string, func(table Blueprint)) error
	Drop(string) error
	Alter(string, func(table Blueprint)) error
	HasTable(string) bool
	MustGet(string) Blueprint
	MustCreate(string, func(table Blueprint)) Blueprint
	MustDrop(string)
	MustAlter(string, func(table Blueprint)) Blueprint
	Rename(string, string) error
	MustRename(string, string) Blueprint
	DropIfExists(string) error
	MustDropIfExists(string)
}

// Blueprint the table operating interface
type Blueprint interface {
	GetName() string
	GetColumns() map[string]*Column
	GetIndexes() map[string]*Index
	GetColumn(name string) *Column
	GetIndex(name string) *Index
	HasColumn(name ...string) bool
	DropColumn(name ...string)
	RenameColumn(old string, new string) *Column
	HasIndex(name ...string) bool
	CreatePrimary(columnName string)
	CreateIndex(key string, columnNames ...string)
	CreateUnique(key string, columnNames ...string)
	DropIndex(key ...string)
	RenameIndex(old string, new string) *Index
	String(name string, length int) *Column
	BigInteger(name string) *Column
	UnsignedBigInteger(name string) *Column
	BigIncrements(name string) *Column
	ID(name string) *Column
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
