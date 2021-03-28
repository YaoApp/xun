package query

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
)

// Builder the dbal query builder
type Builder struct {
	Conn     *Connection
	Attr     Attribute
	Mode     string
	Database string
	Schema   string
	Grammar  dbal.Grammar
}

// Connection DB Connection
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *dbal.Config
	Read        *sqlx.DB
	ReadConfig  *dbal.Config
	Option      *dbal.Option
}

// Table the from attribute
type Table struct {
	Prefix *string
	Name   string
	Alias  string
}

// Where The where constraint for the query.
type Where struct {
	Type     string // where, nested, sub
	Column   string
	Operator string
	Value    interface{}
	Boolean  string
	Wheres   []Where
	Query    *Builder
}

// Attribute the builder Attribute
type Attribute struct {
	From     Table
	Bindings map[string][]interface{}
	Wheres   []Where
}
