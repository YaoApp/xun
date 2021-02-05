package schema

import "github.com/jmoiron/sqlx"

// Connection DB Connection
type Connection struct{ Write *sqlx.DB }

// Schema The database Schema interface
type Schema interface {
	Create()
	Drop()
	DropIfExists()
	Rename()
}

// BlueprintAPI  the bluprint interface
type BlueprintAPI interface {
	BigInteger()
	String()
	Primary()
}

// Builder the dbal schema driver
type Builder struct {
	Conn *Connection
	Schema
}

// Blueprint the table blueprint
type Blueprint struct {
	BlueprintAPI
}
