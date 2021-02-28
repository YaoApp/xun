package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
)

// Connection the database connection for schema operating
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *dbal.Config
	Option      *dbal.Option
}

// Builder the table schema builder struct
type Builder struct {
	Conn *Connection
	dbal.Grammar
	Mode       string
	DBName     string
	SchemaName string
}

// Table the table struct
type Table struct {
	*dbal.Table
	*Builder
	*Primary
	ColumnMap map[string]*Column
	IndexMap  map[string]*Index
}

// Column the table column struct
type Column struct {
	*dbal.Column
	Table *Table
}

// Index the table index struct
type Index struct {
	*dbal.Index
	Table *Table
}

// Primary the table primary key
type Primary struct {
	*dbal.Primary
	Table *Table
}
