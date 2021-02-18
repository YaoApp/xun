package schema

import (
	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // Load postgres driver
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/mysql"
	"github.com/yaoapp/xun/grammar/postgres"
	"github.com/yaoapp/xun/grammar/sqlite3"
	"github.com/yaoapp/xun/utils"
)

// New create a new schema interface using the given driver and DSN
func New(driver string, dsn string) Schema {
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
	return Use(conn)
}

// Use create a new schema interface using the given connection
func Use(conn *Connection) Schema {
	return NewBuilder(conn)
}

// NewBuilder create a new schema builder instance
func NewBuilder(conn *Connection) *Builder {
	return &Builder{
		Conn:    conn,
		Grammar: NewGrammar(conn.WriteConfig.Driver),
		Mode:    "production",
	}
}

// NewGrammar create a new grammar instance
func NewGrammar(driver string) grammar.Grammar {
	switch driver {
	case "mysql":
		return mysql.New()
	case "sqlite3":
		return sqlite3.New()
	case "postgres":
		return postgres.New()
	}
	return mysql.New()
}

// Table create the table blueprint instance
func (builder *Builder) table(name string) *Table {
	table := NewTable(name, builder)
	return table
}

// HasTable determine if the given table exists.
func (builder *Builder) HasTable(name string) bool {
	return builder.Grammar.Exists(name, builder.Conn.Write)
}

// Get a table on the schema.
func (builder *Builder) Get(name string) (Blueprint, error) {
	table := builder.table(name)
	err := builder.Grammar.Get(&table.Table, builder.Conn.Write)
	if err != nil {
		return nil, err
	}

	// attaching columns
	for _, column := range table.Table.Columns {
		name := column.Name
		table.ColumnMap[name] = &Column{
			Column: *column,
			Table:  table,
		}
	}
	// attaching indexes
	for _, index := range table.Table.Indexes {
		name := index.Name
		table.IndexMap[name] = &Index{
			Index: *index,
			Table: table,
		}
	}
	return table, nil
}

// MustGet a table on the schema.
func (builder *Builder) MustGet(name string) Blueprint {
	table, err := builder.Get(name)
	utils.PanicIF(err)
	return table
}

// Create a new table on the schema.
func (builder *Builder) Create(name string, callback func(table Blueprint)) error {
	table := builder.table(name)
	callback(table)
	return builder.Grammar.Create(&table.Table, builder.Conn.Write)
}

// MustCreate a new table on the schema.
func (builder *Builder) MustCreate(name string, callback func(table Blueprint)) Blueprint {
	table := builder.table(name)
	callback(table)
	err := builder.Grammar.Create(&table.Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
}

// Alter a table on the schema.
func (builder *Builder) Alter(name string, callback func(table Blueprint)) error {
	table := builder.table(name)
	callback(table)
	return builder.Grammar.Alter(&table.Table, builder.Conn.Write)
}

// MustAlter a table on the schema.
func (builder *Builder) MustAlter(name string, callback func(table Blueprint)) Blueprint {
	table := builder.table(name)
	callback(table)
	err := builder.Grammar.Alter(&table.Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
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

// DropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropIfExists(name string) error {
	return builder.Grammar.DropIfExists(name, builder.Conn.Write)
}

// MustDropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) MustDropIfExists(name string) {
	err := builder.DropIfExists(name)
	utils.PanicIF(err)
}

// Rename a table on the schema.
func (builder *Builder) Rename(old string, new string) error {
	return builder.Grammar.Rename(old, new, builder.Conn.Write)
}

// MustRename a table on the schema.
func (builder *Builder) MustRename(old string, new string) Blueprint {
	err := builder.Rename(old, new)
	utils.PanicIF(err)
	return builder.table(new)
}

// Primary Specify the primary key(s) for the table.
func (builder *Builder) Primary() {}
