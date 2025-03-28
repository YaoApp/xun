package schema

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// New create a new schema interface using the given driver and DSN
func New(driver string, dsn string) Schema {
	db, err := sqlx.Connect(driver, dsn)
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
	return useBuilder(conn)
}

// newBuilder New create a new schema builder interface using the given driver and DSN
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
			Name:   "main",
		},
		Option: &dbal.Option{},
	}
	return useBuilder(conn)
}

// useBuilder create a new schema builder instance using the given connection
func useBuilder(conn *Connection) *Builder {
	grammar := newGrammar(conn)
	builder := Builder{
		Mode:     "production",
		Conn:     conn,
		Grammar:  grammar,
		Database: grammar.GetDatabase(),
		Schema:   grammar.GetSchema(),
	}
	return &builder
}

// newGrammar create a new grammar interface
func newGrammar(conn *Connection) dbal.Grammar {
	driver := conn.WriteConfig.Driver
	grammar, has := dbal.Grammars[driver]
	if !has {
		panic(fmt.Errorf("The %s driver not import", driver))
	}
	// create ne grammar using the registered grammars
	grammar, err := grammar.NewWith(conn.Write, conn.WriteConfig, conn.Option)
	if err != nil {
		panic(fmt.Errorf("grammar setup error. (%s)", err))
	}
	err = grammar.OnConnected()
	if err != nil {
		panic(fmt.Errorf("the OnConnected event error. (%s)", err))
	}
	return grammar
}

// Builder get the query builder instance
func (builder *Builder) Builder() *Builder {
	return builder
}

// reconnect reconnect db server using setting driver and dsn
func (builder *Builder) reconnect() {
	driver := builder.Conn.WriteConfig.Driver
	dsn := builder.Conn.WriteConfig.DSN
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		panic(err)
	}
	builder.Conn.Write = db
	builder.Grammar, _ = builder.Grammar.NewWith(builder.Conn.Write, builder.Conn.WriteConfig, builder.Conn.Option)
}

// Table create the table blueprint instance
func (builder *Builder) table(name string) *Table {
	table := NewTable(name, builder)
	return table
}

// SetOption set the option of connection
func (builder *Builder) SetOption(option *dbal.Option) {
	builder.Conn.Option = option
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

// GetDB Get the sqlx DB instance
func (builder *Builder) GetDB() (*sqlx.DB, error) {
	if builder.Conn == nil || builder.Conn.Write == nil {
		return nil, fmt.Errorf("the connection is nil")
	}
	return builder.Conn.Write, nil
}

// MustGetDB Get the sqlx DB instance
func (builder *Builder) MustGetDB() *sqlx.DB {
	db, err := builder.GetDB()
	utils.PanicIF(err)
	return db
}

// DB  Alias MustGetDB Get the sqlx DB instance
func (builder *Builder) DB() *sqlx.DB {
	return builder.MustGetDB()
}

// MustGetConnection Get the database connection instance.
func (builder *Builder) MustGetConnection() *dbal.Connection {
	connection, err := builder.GetConnection()
	utils.PanicIF(err)
	return connection
}

// GetTables Get all of the table names for the schema.
func (builder *Builder) GetTables() ([]string, error) {
	tables, err := builder.Grammar.GetTables()
	if err != nil {
		return nil, err
	}

	// - prefix
	if builder.Conn.Option.Prefix != "" {
		for i, tab := range tables {
			tables[i] = strings.TrimLeft(tab, builder.Conn.Option.Prefix)
		}
	}

	return tables, nil
}

// MustGetTables Get all of the table names for the schema.
func (builder *Builder) MustGetTables() []string {
	tables, err := builder.GetTables()
	utils.PanicIF(err)
	return tables
}

// HasTable determine if the given table exists.
func (builder *Builder) HasTable(name string) (bool, error) {
	table := builder.table(name)
	return builder.Grammar.TableExists(table.GetFullName())
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
	dbalTable, err := builder.Grammar.GetTable(table.GetFullName())
	if err != nil {
		return nil, err
	}

	table.Table = dbalTable

	// attaching columns
	for _, column := range table.Table.Columns {
		name := column.Name
		table.ColumnNames = append(table.ColumnNames, name)
		table.ColumnMap[name] = &Column{
			Column: column,
			Table:  table,
		}
	}

	// attaching indexes
	for _, index := range table.Table.Indexes {
		name := index.Name
		table.IndexNames = append(table.IndexNames, name)
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
func (builder *Builder) CreateTable(name string, callback func(table Blueprint), options ...dbal.CreateTableOption) error {
	table := builder.table(name)
	callback(table)
	err := builder.Grammar.CreateTable(table.Table, options...)
	if err != nil {
		return err
	}
	return nil
}

// MustCreateTable create a new table on the schema.
func (builder *Builder) MustCreateTable(name string, callback func(table Blueprint), options ...dbal.CreateTableOption) {
	err := builder.CreateTable(name, callback, options...)
	utils.PanicIF(err)
}

// AlterTable alter a table on the schema.
func (builder *Builder) AlterTable(name string, callback func(table Blueprint)) error {
	table := builder.MustGetTable(name)
	callback(table)
	err := builder.Grammar.AlterTable(table.Get().Table)
	if err != nil {
		return err
	}
	return nil
}

// MustAlterTable alter a table on the schema.
func (builder *Builder) MustAlterTable(name string, callback func(table Blueprint)) {
	err := builder.AlterTable(name, callback)
	utils.PanicIF(err)
}

// DropTable Indicate that the table should be dropped.
func (builder *Builder) DropTable(name string) error {
	table := builder.table(name)
	return builder.Grammar.DropTable(table.GetFullName())
}

// MustDropTable Indicate that the table should be dropped.
func (builder *Builder) MustDropTable(name string) {
	table := builder.table(name)
	err := builder.DropTable(table.GetFullName())
	utils.PanicIF(err)
}

// DropTableIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) DropTableIfExists(name string) error {
	table := builder.table(name)
	return builder.Grammar.DropTableIfExists(table.GetFullName())
}

// MustDropTableIfExists Indicate that the table should be dropped if it exists.
func (builder *Builder) MustDropTableIfExists(name string) {
	err := builder.DropTableIfExists(name)
	utils.PanicIF(err)
}

// RenameTable rename a table on the schema.
func (builder *Builder) RenameTable(old string, new string) error {
	oldTab := builder.table(old)
	newTab := builder.table(new)
	return builder.Grammar.RenameTable(oldTab.GetFullName(), newTab.GetFullName())
}

// MustRenameTable rename a table on the schema.
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
	version, err := builder.Grammar.GetVersion()
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
