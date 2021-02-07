package query

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal/schema"
)

// Query The database Query interface
type Query interface {
	Where()
	Join()
}

// Builder the dbal query builder
type Builder struct {
	Conn *Connection
	Query
}

// Connection DB Connection
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *schema.ConnConfig
	Read        *sqlx.DB
	ReadConfig  *schema.ConnConfig
	Config      *schema.Config
}
