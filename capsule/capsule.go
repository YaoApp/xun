package capsule

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/model"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Global The global manager
var Global *Manager = nil

// New Create a database manager instance.
func New() *Manager {
	return &Manager{
		Pool:        &Pool{},
		Connections: &sync.Map{},
		Option:      &dbal.Option{},
	}
}

// NewWithOption Create a database manager instance using the given option.
func NewWithOption(option dbal.Option) *Manager {
	manager := New()
	manager.SetOption(option)
	return manager
}

// AddConn Register a connection with the manager.
func AddConn(name string, driver string, datasource string) *Manager {
	return New().AddConn(name, driver, datasource)
}

// AddConn Register a connection with the manager.
func (manager *Manager) AddConn(name string, driver string, datasource string) *Manager {
	manager.AddConnection(name, driver, datasource, false)
	return manager
}

// AddReadConn Register a readonly connection with the manager.
func AddReadConn(name string, driver string, datasource string) *Manager {
	return New().AddReadConn(name, driver, datasource)
}

// AddReadConn Register a readonly with the manager.
func (manager *Manager) AddReadConn(name string, driver string, datasource string) *Manager {
	manager.AddConnection(name, driver, datasource, true)
	return manager
}

// SetOption set the database manager as the given value
func (manager *Manager) SetOption(option dbal.Option) {
	manager.Option = &option
}

// AddConnection Register a connection with the manager.
func (manager *Manager) AddConnection(name string, driver string, datasource string, readonly bool) *Manager {
	config := dbal.Config{
		Name:     name,
		Driver:   driver,
		DSN:      datasource,
		ReadOnly: readonly,
	}

	conn := &Connection{
		DB:     *sqlx.MustOpen(config.Driver, config.DSN),
		Config: &config,
	}
	manager.Pool.Primary = append(manager.Pool.Primary, conn)
	if config.ReadOnly == true {
		manager.Pool.Readonly = append(manager.Pool.Readonly, conn)
	} else {
		manager.Pool.Primary = append(manager.Pool.Primary, conn)
	}
	manager.Connections.Store(config.Name, conn)

	if Global == nil {
		Global = manager
	}

	return manager
}

// GetConnection Get a registered connection instance.
func (manager *Manager) GetConnection(name string) *Connection {

	c, has := manager.Connections.Load(name)
	conn := c.(*Connection)
	if !has {
		err := errors.New("the connection " + name + " is not registered")
		panic(err)
	}
	return conn
}

// GetRand Get a registered connection instance.
func GetRand(connections []*Connection) *Connection {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	i := r.Intn(len(connections))
	return connections[i]
}

// GetPrimary Get a registered primary connection instance.
func (manager *Manager) GetPrimary() *Connection {
	length := len(manager.Pool.Primary)
	if length < 1 {
		err := errors.New("the Primary connection not found ")
		panic(err)
	} else if length == 1 {
		return manager.Pool.Primary[0]
	}
	return GetRand(manager.Pool.Primary)
}

// GetRead Get a registered read only connection instance.
func (manager *Manager) GetRead() *Connection {
	length := len(manager.Pool.Readonly)
	if length < 1 {
		return manager.GetPrimary()
	} else if length == 1 {
		return manager.Pool.Readonly[0]
	}
	return GetRand(manager.Pool.Readonly)
}

// SetAsGlobal Make this connetion instance available globally.
func (manager *Manager) SetAsGlobal() {
	Global = manager
}

// Schema Get a schema builder instance.
func Schema() schema.Schema {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Schema()
}

// Schema Get a schema builder instance.
func (manager *Manager) Schema() schema.Schema {
	write := manager.GetPrimary()
	return schema.Use(&schema.Connection{
		Write:       &write.DB,
		WriteConfig: write.Config,
		Option:      manager.Option,
	})
}

// Query Get a fluent query builder instance.
func Query() query.Query {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Query()
}

// Query Get a fluent query builder instance.
func (manager *Manager) Query() query.Query {
	write := manager.GetPrimary()
	read := manager.GetRead()
	return query.Use(
		&query.Connection{
			Write:       &write.DB,
			WriteConfig: write.Config,
			Read:        &read.DB,
			ReadConfig:  read.Config,
			Option:      manager.Option,
		})
}

// Build to build a new xun model instance
func Build(v interface{}, flow ...interface{}) *model.Model {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Build(v, flow...)
}

// Build to build a new xun model instance
func (manager *Manager) Build(v interface{}, flow ...interface{}) *model.Model {
	query := manager.Query()
	schema := manager.Schema()
	return model.Build(query.Builder(), schema, v, flow...)
}

// Make to create a new xun registered model instance
func Make(name string) interface{} {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Make(name)
}

// Make to create a new xun model instance
func (manager *Manager) Make(name string) interface{} {
	query := manager.Query()
	schema := manager.Schema()
	return model.New(name, query.Builder(), schema)
}
