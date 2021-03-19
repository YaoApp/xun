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

// NewHook get the hook of driver( should be optimized )
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

// GetConnection Get the database connection instance.
func (builder *Builder) GetConnection() (*dbal.Connection, error) {
	version, err := builder.GetVersion()
	if err != nil {
		return nil, err
	}
	return &dbal.Connection{
		DB:      builder.Conn.Write,
		Config:  builder.Conn.WriteConfig,
		Option:  builder.Conn.Option,
		Version: version,
	}, nil
}

// MustGetConnection Get the database connection instance.
func (builder *Builder) MustGetConnection() *dbal.Connection {
	connection, err := builder.GetConnection()
	utils.PanicIF(err)
	return connection
}

// GetTables Get all of the table names for the schema.
func (builder *Builder) GetTables() ([]string, error) {
	return builder.Grammar.GetTables(builder.Conn.Write)
}

// MustGetTables Get all of the table names for the schema.
func (builder *Builder) MustGetTables() []string {
	tables, err := builder.GetTables()
	utils.PanicIF(err)
	return tables
}

// HasTable determine if the given table exists.
func (builder *Builder) HasTable(name string) (bool, error) {
	return builder.Grammar.TableExists(name, builder.Conn.Write)
}

// MustHasTable determine if the given table exists.
func (builder *Builder) MustHasTable(name string) bool {
	has, err := builder.HasTable(name)
	utils.PanicIF(err)
	return has
}

// GetTable a table on the schema.
func (builder *Builder) GetTable(name string) (Blueprint, error) {
	table := builder.table(name)
	err := builder.Grammar.GetTable(table.Table, builder.Conn.Write)
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

// MustGetTable a table on the schema.
func (builder *Builder) MustGetTable(name string) Blueprint {
	table, err := builder.GetTable(name)
	utils.PanicIF(err)
	return table
}

// CreateTable create a new table on the schema.
func (builder *Builder) CreateTable(name string, callback func(table Blueprint)) error {
	table := builder.table(name)
	callback(table)
	return builder.Grammar.CreateTable(table.Table, builder.Conn.Write)
}

// MustCreateTable create a new table on the schema.
func (builder *Builder) MustCreateTable(name string, callback func(table Blueprint)) Blueprint {
	table := builder.table(name)
	callback(table)
	err := builder.Grammar.CreateTable(table.Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
}

// AlterTable alter a table on the schema.
func (builder *Builder) AlterTable(name string, callback func(table Blueprint)) error {
	table := builder.MustGetTable(name)
	callback(table)
	return builder.Grammar.AlterTable(table.Get().Table, builder.Conn.Write)
}

// MustAlterTable alter a table on the schema.
func (builder *Builder) MustAlterTable(name string, callback func(table Blueprint)) Blueprint {
	table := builder.MustGetTable(name)
	callback(table)
	err := builder.Grammar.AlterTable(table.Get().Table, builder.Conn.Write)
	utils.PanicIF(err)
	return table
}

// DropTable Indicate that the table should be dropped.
func (builder *Builder) DropTable(name string) error {
	return builder.Grammar.DropTable(name, builder.Conn.Write)
}

// MustDropTable Indicate that the table should be dropped.
func (builder *Builder) MustDropTable(name string) {
	err := builder.DropTable(name)
	utils.PanicIF(err)
}

// DropTableIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropTableIfExists(name string) error {
	return builder.Grammar.DropTableIfExists(name, builder.Conn.Write)
}

// MustDropTableIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) MustDropTableIfExists(name string) {
	err := builder.DropTableIfExists(name)
	utils.PanicIF(err)
}

// RenameTable rename a table on the schema.
func (builder *Builder) RenameTable(old string, new string) error {
	return builder.Grammar.RenameTable(old, new, builder.Conn.Write)
}

//MustRenameTable rename a table on the schema.
func (builder *Builder) MustRenameTable(old string, new string) Blueprint {
	err := builder.RenameTable(old, new)
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
