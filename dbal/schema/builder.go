package schema

import (
	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/mysql"
	"github.com/yaoapp/xun/grammar/sqlite3"
	"github.com/yaoapp/xun/utils"
)

// New create new schema buider interface
func New(conn *Connection) Schema {
	builder := NewBuilder(conn)
	return &builder
}

// NewBuilder create a new schema buider blueprint
func NewBuilder(conn *Connection) Builder {
	return Builder{
		Conn:    conn,
		Grammar: NewGrammar(conn.WriteConfig.Driver),
	}
}

// NewGrammar create a new grammar intance
func NewGrammar(driver string) grammar.Grammar {
	switch driver {
	case "mysql":
		return mysql.New()
	case "sqlite3":
		return sqlite3.New()
	}
	return mysql.New()
}

// NewBuilderByDSN create a new schema builder by given DSN
func NewBuilderByDSN(driver string, dsn string) *Builder {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	conn := &Connection{
		Write: db,
		WriteConfig: &dbal.Config{
			DSN:    dsn,
			Driver: driver,
			Name:   "main",
		},
		Config: &dbal.DBConfig{},
	}
	builder := NewBuilder(conn)
	return &builder
}

// Create a new table on the schema.
func (builder *Builder) Create(name string, callback func(table *Blueprint)) error {
	table := builder.Table(name)
	callback(table)
	table.Table = table.GrammarTable()
	return builder.Grammar.Create(table.Table, builder.Conn.Write)
}

// MustCreate a new table on the schema.
func (builder *Builder) MustCreate(name string, callback func(table *Blueprint)) *Blueprint {
	table := builder.Table(name)
	callback(table)
	table.Table = table.GrammarTable()
	err := builder.Grammar.Create(table.Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
}

// Alter a table on the schema.
func (builder *Builder) Alter(name string, callback func(table *Blueprint)) error {
	table := builder.Table(name)
	return table.Alter(callback)
}

// Drop Indicate that the table should be dropped.
func (builder *Builder) Drop(name string) error {
	return builder.Grammar.Drop(name, builder.Conn.Write)
}

// MustDrop Indicate that the table should be dropped.
func (builder *Builder) MustDrop(name string) {
	err := builder.Drop(name)
	utils.PanicIF(err)
}

// Table get the table blueprint instance
func (builder *Builder) Table(name string) *Blueprint {
	table := NewBlueprint(name, builder)
	return table
}

// HasTable determine if the given table exists.
func (builder *Builder) HasTable(name string) bool {
	return builder.Grammar.Exists(name, builder.Conn.Write)
}

// DropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropIfExists(name string) error {
	table := builder.Table(name)
	return table.DropIfExists()
}

// MustDropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) MustDropIfExists(name string) {
	table := builder.Table(name)
	table.MustDropIfExists()
}

// Rename a table on the schema.
func (builder *Builder) Rename(from string, to string) error {
	table := builder.Table(from)
	return table.Rename(to)
}

// MustRename a table on the schema.
func (builder *Builder) MustRename(from string, to string) *Blueprint {
	table := builder.Table(from)
	return table.MustRename(to)
}

// Primary Specify the primary key(s) for the table.
func (builder *Builder) Primary() {}
