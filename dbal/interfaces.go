package dbal

import "github.com/jmoiron/sqlx"

// Grammar the SQL Grammar inteface
type Grammar interface {
	Setup(db *sqlx.DB, config *Config, option *Option) error
	OnConnected() error

	GetVersion(db *sqlx.DB) (*Version, error)
	GetDatabase() string
	GetSchema() string

	GetTables(db *sqlx.DB) ([]string, error)

	TableExists(name string, db *sqlx.DB) (bool, error)
	GetTable(table *Table, db *sqlx.DB) error
	CreateTable(table *Table, db *sqlx.DB) error
	AlterTable(table *Table, db *sqlx.DB) error
	DropTable(name string, db *sqlx.DB) error
	DropTableIfExists(name string, db *sqlx.DB) error
	RenameTable(old string, new string, db *sqlx.DB) error
	GetColumnListing(dbName string, tableName string, db *sqlx.DB) ([]*Column, error)
}

// Quoter the database quoting query text intrface
type Quoter interface {
	ID(name string, db *sqlx.DB) string
	VAL(v interface{}, db *sqlx.DB) string // operates on both string and []byte and int or other types.
}
