package schema

import (
	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
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
	}
	builder := NewBuilder(conn)
	return &builder
}

// Create a new table on the schema.
func (builder *Builder) Create(name string, callback func(table *Blueprint)) {
	table := NewBlueprint(name, builder)
	callback(table)
	table.Create()
}

// Drop Indicate that the table should be dropped.
func (builder *Builder) Drop(name string) {
	table := NewBlueprint(name, builder)
	table.Drop()
}

// Table get the table blueprint instance
func (builder *Builder) Table(name string) *Blueprint {
	table := NewBlueprint(name, builder)
	return table
}

// DropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropIfExists() {}

// Rename the table to a given name.
func (builder *Builder) Rename() {}

// Primary Specify the primary key(s) for the table.
func (builder *Builder) Primary() {}
