package mysql

import (
	"github.com/go-sql-driver/mysql"
)

// DBName Get the database name of the connection
func (grammarSQL *MySQL) DBName() string {

	cfg, err := mysql.ParseDSN(grammarSQL.DSN)
	if err != nil {
		panic(err)
	}

	grammarSQL.DB = cfg.DBName
	return grammarSQL.DB
}

// SchemaName Get the schema name of the connection
func (grammarSQL *MySQL) SchemaName() string {
	schema := grammarSQL.DBName()
	grammarSQL.Schema = schema
	return grammarSQL.Schema
}
