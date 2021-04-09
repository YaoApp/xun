package dbal

import (
	"github.com/jmoiron/sqlx"
)

// Grammar the SQL Grammar inteface
type Grammar interface {
	NewWith(db *sqlx.DB, config *Config, option *Option) (Grammar, error)
	NewWithRead(write *sqlx.DB, writeConfig *Config, read *sqlx.DB, readConfig *Config, option *Option) (Grammar, error)

	OnConnected() error

	GetVersion() (*Version, error)
	GetDatabase() string
	GetSchema() string
	GetOperators() []string
	GetSelectComponents() []string

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
	CompileInsert(query *Query, columns []interface{}, values [][]interface{}) (string, []interface{})
	CompileInsertIgnore(query *Query, columns []interface{}, values [][]interface{}) (string, []interface{})
	CompileInsertGetID(query *Query, columns []interface{}, values [][]interface{}, sequence string) (string, []interface{})
	CompileInsertUsing(query *Query, columns []interface{}, sql string) string
	CompileUpsert(query *Query, columns []interface{}, values [][]interface{}, uniqueBy []interface{}, updateValues interface{}) (string, []interface{})
	CompileSelect(query *Query) string
	CompileSelectOffset(query *Query, offset *int) string

	ProcessInsertGetID(sql string, bindings []interface{}, sequence string) (int64, error)
}

// Quoter the database quoting query text intrface
type Quoter interface {
	Bind(db *sqlx.DB, prefix string, dbRead ...*sqlx.DB) Quoter
	ID(value string) string
	VAL(value interface{}) string // operates on both string and []byte and int or other types.
	Wrap(value interface{}) string
	WrapTable(value interface{}) string
	WrapUnion(sql string) string
	IsExpression(value interface{}) bool
	Parameter(value interface{}, num int) string
	Parameterize(values []interface{}, offset int) string
	Columnize(columns []interface{}) string
}
