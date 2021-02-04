package query

import (
	"github.com/yaoapp/xun/query/driver/mysql"
	"github.com/yaoapp/xun/query/driver/oracle"
	"github.com/yaoapp/xun/query/driver/sqlite"
	"github.com/yaoapp/xun/query/driver/sqlserver"
)

// Table Get a fluent query builder instance.
func Table() Query {
	return mysql.Table()
}

// TableMySQL  Get a fluent query builder instance of MySQL.
func TableMySQL() Query {
	return mysql.Table()
}

// TableSQLite Get a fluent query builder instance of SQLite.
func TableSQLite() Query {
	return sqlite.Table()
}

// TableSQLServer Get a fluent query builder instance of SQL Server.
func TableSQLServer() Query {
	return sqlserver.Table()
}

// TableOracle Get a fluent query builder instance of Oracle.
func TableOracle() Query {
	return oracle.Table()
}
