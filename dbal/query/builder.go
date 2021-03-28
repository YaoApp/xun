package query

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"

	_ "github.com/yaoapp/xun/grammar/mysql"    // Load the MySQL Grammar
	_ "github.com/yaoapp/xun/grammar/postgres" // Load the Postgres Grammar
	_ "github.com/yaoapp/xun/grammar/sqlite3"  // Load the SQLite3 Grammar
)

// New create a new schema interface using the given driver and DSN
func New(driver string, dsn string) Query {
	builder := newBuilder(driver, dsn)
	return builder
}

// Use create a new schema interface using the given connection
func Use(conn *Connection) Query {
	builder := useBuilder(conn)
	return builder
}

// DB Get the sqlx.DB pointer instance
func (builder *Builder) DB(readonly ...bool) *sqlx.DB {
	if len(readonly) == 1 && readonly[0] == true {
		return builder.Conn.Read
	}
	return builder.Conn.Write
}

// clone create a new builder instance with current builder
func (builder *Builder) clone() *Builder {
	new := *builder
	return &new
}

// new create a new builder instance
func (builder *Builder) new() *Builder {
	new := *builder
	new.renewAttribute()
	return &new
}

// newBuilder create a new schema builder interface using the given driver and DSN
func newBuilder(driver string, dsn string) *Builder {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		panic(err)
	}
	conn := &Connection{
		Write: db,
		WriteConfig: &dbal.Config{
			DSN:    dsn,
			Driver: driver,
			Name:   "primary",
		},
		Read: db,
		ReadConfig: &dbal.Config{
			DSN:      dsn,
			Driver:   driver,
			Name:     "secondary",
			ReadOnly: true,
		},
		Option: &dbal.Option{},
	}
	return useBuilder(conn)
}

// useBuilder create a new schema builder instance using the given connection
func useBuilder(conn *Connection) *Builder {
	grammar := newGrammar(conn)
	return &Builder{
		Mode:     "production",
		Conn:     conn,
		Grammar:  grammar,
		Database: grammar.GetDatabase(),
		Schema:   grammar.GetSchema(),
		Attr:     newAttribute(),
	}
}

// newGrammar create a new grammar interface
func newGrammar(conn *Connection) dbal.Grammar {
	driver := conn.WriteConfig.Driver
	grammar, has := dbal.Grammars[driver]
	if !has {
		panic(fmt.Errorf("The %s driver not import", driver))
	}
	// create ne grammar using the registered grammars
	grammar, err := grammar.NewWithRead(conn.Write, conn.WriteConfig, conn.Read, conn.ReadConfig, conn.Option)
	if err != nil {
		panic(fmt.Errorf("grammar setup error. (%s)", err))
	}
	err = grammar.OnConnected()
	if err != nil {
		panic(fmt.Errorf("the OnConnected event error. (%s)", err))
	}
	return grammar
}

// newAttribute create a new Attribute instance
func newAttribute() Attribute {
	return Attribute{}
}

// renewAttribute reset and create a new Attribute instance
func (builder *Builder) renewAttribute() {
	builder.Attr = newAttribute()
}
