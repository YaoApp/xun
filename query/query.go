package query

import (
	dbal "github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/query/driver/mysql"
	"github.com/yaoapp/xun/query/driver/oracle"
	"github.com/yaoapp/xun/query/driver/postgresql"
	"github.com/yaoapp/xun/query/driver/sqlite"
	"github.com/yaoapp/xun/query/driver/sqlserver"
)

// New Get a fluent query builder instance.
func New(driver string, conn *dbal.Connection) Query {
	switch driver {
	case "mysql":
		return mysql.New(conn)
	case "SQLite":
		return sqlite.New(conn)
	case "sqlsvr":
		return sqlserver.New(conn)
	case "oracle":
		return oracle.New(conn)
	case "postgresql":
		return postgresql.New(conn)
	default:
		return mysql.New(conn)
	}
}
