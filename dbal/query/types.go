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

// Name the from attribute
type Name struct {
	Prefix *string
	Name   string
	Alias  string
}

// Where The where constraint for the query.
type Where struct {
	Type     string // basic, nested, sub
	Column   string
	Operator string
	Value    interface{}
	Boolean  string
	Wheres   []Where
	Query    *Builder
}

// Attribute the builder Attribute
type Attribute struct {
	From     Name
	Bindings map[string][]interface{}
	Wheres   []Where
	Columns  []Name
}
