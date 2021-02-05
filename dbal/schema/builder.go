package schema

import (
	"fmt"

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
	table := NewBlueprint(name)
	callback(table)
	sql := table.sqlCreate()
	builder.Conn.Write.MustExec(sql)
}

// Drop Indicate that the table should be dropped.
func (builder *Builder) Drop() {
	fmt.Printf("\nDrop DBAL: \n===\n%#v\n===\n", builder.Conn)
}

// DropIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropIfExists() {}

// Rename the table to a given name.
func (builder *Builder) Rename() {}

// Primary Specify the primary key(s) for the table.
func (builder *Builder) Primary() {}
