package dbal

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun"
)

// Grammar the SQL Grammar inteface
type Grammar interface {
	NewWith(db *sqlx.DB, config *Config, option *Option) (Grammar, error)
	NewWithRead(write *sqlx.DB, writeConfig *Config, read *sqlx.DB, readConfig *Config, option *Option) (Grammar, error)

	OnConnected() error

	GetVersion() (*Version, error)
	GetDatabase() string
	GetSchema() string

	// Grammar for migrating
	GetTables() ([]string, error)

	TableExists(name string) (bool, error)
	GetTable(name string) (*Table, error)
	CreateTable(table *Table) error
	AlterTable(table *Table) error
	DropTable(name string) error
	DropTableIfExists(name string) error
	RenameTable(old string, new string) error
	GetColumnListing(dbName string, tableName string) ([]*Column, error)

	// Grammar for querying
	Insert(tableName string, values []xun.R) (sql.Result, error)
	InsertIgnore(tableName string, values []xun.R) (sql.Result, error)
}

// Quoter the database quoting query text intrface
type Quoter interface {
	ID(name string, db *sqlx.DB) string
	VAL(v interface{}, db *sqlx.DB) string // operates on both string and []byte and int or other types.
}
