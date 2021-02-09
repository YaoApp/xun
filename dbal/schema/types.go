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
	MustAlter(string, func(table *Blueprint)) *Blueprint
}

// BlueprintMethods  the bluprint interface
type BlueprintMethods interface {
	BigInteger()
	String(name string, length int) *Blueprint
	Primary()
}

// Builder the dbal schema driver
type Builder struct {
	Conn    *Connection
	Grammar grammar.Grammar
	Schema
}

// Blueprint the table blueprint
type Blueprint struct {
	BlueprintMethods
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
	Table    *Blueprint
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
	Table   *Blueprint
	dropped bool
	renamed bool
	newname string
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
