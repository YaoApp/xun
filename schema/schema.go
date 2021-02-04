package schema

import (
	dbal "github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/schema/driver/mysql"
	"github.com/yaoapp/xun/schema/driver/oracle"
	"github.com/yaoapp/xun/schema/driver/postgresql"
	"github.com/yaoapp/xun/schema/driver/sqlite"
	"github.com/yaoapp/xun/schema/driver/sqlserver"
)

// New create the schema instance()
func New(driver string, conn *dbal.Connection) Schema {
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
