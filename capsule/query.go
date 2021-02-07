package capsule

import (
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/query/mysql"
	"github.com/yaoapp/xun/query/oracle"
	"github.com/yaoapp/xun/query/postgresql"
	"github.com/yaoapp/xun/query/sqlite3"
	"github.com/yaoapp/xun/query/sqlserver"
)

// newQuery Get a query builder instance.
func newQuery(driver string, conn *query.Connection) query.Query {
	switch driver {
	case "mysql":
		return mysql.New(conn)
	case "sqlite3":
		return sqlite3.New(conn)
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
