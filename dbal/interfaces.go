package dbal

import "github.com/jmoiron/sqlx"

// Grammar the SQL Grammar inteface
type Grammar interface {
	NewWith(db *sqlx.DB, config *Config, option *Option) (Grammar, error)
	OnConnected() error

	GetVersion() (*Version, error)
	GetDatabase() string
	GetSchema() string

	GetTables() ([]string, error)

	TableExists(name string) (bool, error)
	GetTable(table *Table) error
	CreateTable(table *Table) error
	AlterTable(table *Table) error
	DropTable(name string) error
	DropTableIfExists(name string) error
	RenameTable(old string, new string) error
	GetColumnListing(dbName string, tableName string) ([]*Column, error)
}

// Quoter the database quoting query text intrface
type Quoter interface {
	ID(name string, db *sqlx.DB) string
	VAL(v interface{}, db *sqlx.DB) string // operates on both string and []byte and int or other types.
}
