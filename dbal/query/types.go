package query

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
)

// Builder the dbal query builder
type Builder struct {
	Conn     *Connection
	Query    *dbal.Query
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
