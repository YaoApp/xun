package schema

import "fmt"

// New create new schema buider interface
func New(conn *Connection) Schema {
	blueprint := NewBuilder(conn)
	return &blueprint
}

// NewBuilder create new schema buider blueprint
func NewBuilder(conn *Connection) Builder {
	return Builder{
		Conn: conn,
	}
}

// Create a new table on the schema.
func (builder *Builder) Create() {}

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
