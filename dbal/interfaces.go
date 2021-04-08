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
	Insert(query *Query, values []xun.R) (sql.Result, error)
	InsertIgnore(query *Query, values []xun.R) (sql.Result, error)
	InsertUsing(query *Query, columns []interface{}, sql string, bindings []interface{}) (sql.Result, error)
	InsertGetID(query *Query, values []xun.R, sequence string) (int64, error)

	Upsert(query *Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (sql.Result, error)

	CompileInsert(query *Query, values []xun.R) (string, []interface{})
	CompileUpsert(query *Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (string, []interface{})
	CompileSelect(query *Query) string
	CompileSelectOffset(query *Query, offset *int) string
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
