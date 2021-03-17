package schema

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
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
		Option: &dbal.Option{},
	}
	return Use(conn)
}

// Use create a new schema interface using the given connection
func Use(conn *Connection) Schema {
	return NewBuilder(conn)
}

// NewBuilder create a new schema builder instance
func NewBuilder(conn *Connection) *Builder {
	grammar := NewGrammar(conn.WriteConfig.Driver, conn.WriteConfig.DSN)
	hook := NewHook(conn.WriteConfig.Driver)
	builder := Builder{
		Mode:       "production",
		Conn:       conn,
		Hook:       hook,
		Grammar:    grammar,
		DBName:     grammar.GetDBName(),
		SchemaName: grammar.GetSchemaName(),
	}

	if builder.Hook.OnConnected != nil {
		builder.Hook.OnConnected(grammar, conn.Write)
	}
	return &builder
}

// NewGrammar create a new grammar instance
func NewGrammar(driver string, dsn string) dbal.Grammar {
	grammar, has := dbal.Grammars[driver]
	if !has {
		panic(errors.New("the " + driver + "driver not import!"))
	}
	grammar.Config(dsn)
	return grammar
}

// NewHook get the hook of driver
func NewHook(driver string) dbal.Hook {
	hook, has := dbal.Hooks[driver]
	if !has {
		return dbal.Hook{}
	}
	return hook
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
	err := builder.Grammar.Get(table.Table, builder.Conn.Write)
	if err != nil {
		return nil, err
	}

	// attaching columns
	for _, column := range table.Table.Columns {
		name := column.Name
		table.ColumnMap[name] = &Column{
			Column: column,
			Table:  table,
		}
	}

	// attaching indexes
	for _, index := range table.Table.Indexes {
		name := index.Name
		table.IndexMap[name] = &Index{
			Index: index,
			Table: table,
		}
	}

	// attaching primary
	if table.Table.Primary != nil {
		table.Primary = &Primary{
			Primary: table.Table.Primary,
			Table:   table,
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
	return builder.Grammar.Create(table.Table, builder.Conn.Write)
}

// MustCreate a new table on the schema.
func (builder *Builder) MustCreate(name string, callback func(table Blueprint)) Blueprint {
	table := builder.table(name)
	callback(table)
	err := builder.Grammar.Create(table.Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
}

// Alter a table on the schema.
func (builder *Builder) Alter(name string, callback func(table Blueprint)) error {
	table := builder.MustGet(name)
	callback(table)
	return builder.Grammar.Alter(table.GetTable().Table, builder.Conn.Write)
}

// MustAlter a table on the schema.
func (builder *Builder) MustAlter(name string, callback func(table Blueprint)) Blueprint {
	table := builder.MustGet(name)
	callback(table)
	err := builder.Grammar.Alter(table.GetTable().Table, builder.Conn.Write)
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

// GetVersion get the version of the connection database
func (builder *Builder) GetVersion() (*dbal.Version, error) {

	if builder.Conn.Version != nil {
		return builder.Conn.Version, nil
	}

	// Query Version using connection
	version, err := builder.Grammar.GetVersion(builder.Conn.Write)
	if err != nil {
		return nil, err
	}
	builder.Conn.Version = version
	return version, nil
}

// MustGetVersion get the version of the connection database
func (builder *Builder) MustGetVersion() *dbal.Version {
	version, err := builder.GetVersion()
	utils.PanicIF(err)
	return version
}
