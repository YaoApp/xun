package query

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal/query"
)

// Builder a query builder
type Builder struct{ query.Builder }

// Query The database Query interface
type Query interface{ query.Query }

// Connection DB Connection
type Connection struct {
	Write *sqlx.DB
	Read  *sqlx.DB
}
