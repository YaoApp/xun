package schema

import "github.com/jmoiron/sqlx"

// Schema The database Schema interface
type Schema interface {
	Create()
	Drop()
	DropIfExists()
	Rename()
	Primary()

	BigInteger()
	String()
}

// Blueprint the dbal schema driver
type Blueprint struct {
	Conn *Connection
	Schema
}

// Connection DB Connection
type Connection struct{ Write *sqlx.DB }
