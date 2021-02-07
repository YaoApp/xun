package dbal

import (
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Config the Connection configuration
type Config struct {
	Driver   string `json:"driver"`        // The driver name. mysql,oci,pgsql,sqlsrv,sqlite
	DSN      string `json:"dsn,omitempty"` // The driver wrapper. sqlite:///:memory:, mysql://localhost:4486/foo?charset=UTF8
	Name     string `json:"name,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`
}

// DBConfig the database configuration
type DBConfig struct {
	TablePrefix string `json:"table_prefix,omitempty"`
	DBPrefix    string `json:"db_prefix,omitempty"`
	DBName      string `json:"db,omitempty"`
	Collation   string `json:"collation,omitempty"`
	Charset     string `json:"charset,omitempty"`
}

// DBName parse the DSN, and return the database name
func (config Config) DBName() string {
	switch config.Driver {
	case "mysql":
		cfg, err := mysql.ParseDSN(config.DSN)
		if err != nil {
			panic(err)
		}
		return cfg.DBName
	case "sqlite3":
		return config.Sqlite3DBName()
	}
	return ""
}

// Sqlite3DBName parse the DSN, and return the database name for sqlite3 driver.
func (config Config) Sqlite3DBName() string {
	pos := strings.IndexRune(config.DSN, '?')
	if pos >= 1 {
		params, err := url.ParseQuery(config.DSN[pos+1:])
		if err != nil {
			panic(err)
		}
		return params.Get("url")
	}
	return ""
}
