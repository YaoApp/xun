package query

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
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
	WriteConfig *dbal.Config
	Read        *sqlx.DB
	ReadConfig  *dbal.Config
	Config      *dbal.DBConfig
}
