package schema

import (
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/schema/driver/mysql"
	"github.com/yaoapp/xun/schema/driver/oracle"
	"github.com/yaoapp/xun/schema/driver/sqlite"
	"github.com/yaoapp/xun/schema/driver/sqlserver"
)

// New create the schema instance()
func New() dbal.Schema {
	return mysql.New()
}

// NewMySQL create the MySQL schema instance()
func NewMySQL() dbal.Schema {
	return mysql.New()
}

// NewSQLite create the SQLite schema instance()
func NewSQLite() dbal.Schema {
	return sqlite.New()
}

// NewSQLServer create the SQL Server schema instance()
func NewSQLServer() dbal.Schema {
	return sqlserver.New()
}

// NewOracle create the Oracle schema instance()
func NewOracle() dbal.Schema {
	return oracle.New()
}
