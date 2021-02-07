package sqlite3

import (
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/schema"
)

// New create new mysql blueprint instance
func New(conn *schema.Connection) schema.Schema {
	return &Builder{
		Builder: schema.NewBuilder(conn),
	}
}

// NewBuilder create new schema buider blueprint
func NewBuilder(conn *schema.Connection) Builder {
	return Builder{
		Builder: schema.Builder{
			Conn: conn,
		},
	}
}

// NewBuilderByDSN create a new schema builder by given DSN
func NewBuilderByDSN(driver string, dsn string) *Builder {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	conn := &schema.Connection{
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
