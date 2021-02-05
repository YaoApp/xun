package capsule

import (
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/schema/mysql"
	"github.com/yaoapp/xun/schema/oracle"
	"github.com/yaoapp/xun/schema/postgresql"
	"github.com/yaoapp/xun/schema/sqlite"
	"github.com/yaoapp/xun/schema/sqlserver"
)

// Get a schema builder instance.
func newSchema(driver string, conn *schema.Connection) schema.Schema {
	switch driver {
	case "mysql":
		return mysql.New(conn)
	case "sqlite":
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
