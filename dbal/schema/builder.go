package schema

import (
	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
)

// New create new schema buider interface
func New(conn *Connection) Schema {
	builder := NewBuilder(conn)
	return &builder
}

// NewBuilder create new schema buider blueprint
func NewBuilder(conn *Connection) Builder {
	return Builder{
		Conn: conn,
	}
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
	return table.Create(callback)
}

// MustCreate a new table on the schema.
func (builder *Builder) MustCreate(name string, callback func(table *Blueprint)) *Blueprint {
	table := builder.Table(name)
	return table.MustCreate(callback)
}

// Alter a table on the schema.
func (builder *Builder) Alter(name string, callback func(table *Blueprint)) error {
	table := builder.Table(name)
	return table.Alter(callback)
}

// Drop Indicate that the table should be dropped.
func (builder *Builder) Drop(name string) error {
	table := builder.Table(name)
	return table.Drop()
}

// MustDrop Indicate that the table should be dropped.
func (builder *Builder) MustDrop(name string) {
	table := builder.Table(name)
	table.MustDrop()
}

// Table get the table blueprint instance
func (builder *Builder) Table(name string) *Blueprint {
	table := NewBlueprint(name, builder)
	return table
}

// HasTable determine if the given table exists.
func (builder *Builder) HasTable(name string) bool {
	table := builder.Table(name)
	return table.Exists()
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
