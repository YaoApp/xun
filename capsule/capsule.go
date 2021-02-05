package capsule

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
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
	}
}

// AddConnection Register a connection with the manager.
func AddConnection(config Config) *Manager {
	return New().AddConnection(config)
}

// AddConnection Register a connection with the manager.
func (manager *Manager) AddConnection(config Config) *Manager {

	name := "main"
	if config.Name != "" {
		name = config.Name
	}

	conn := &Connection{
		DB:     *sqlx.MustOpen(config.DriverName(), config.DataSource()),
		Config: &config,
	}

	if config.ReadOnly == true {
		manager.Pool.Readonly = append(manager.Pool.Readonly, conn)
	} else {
		manager.Pool.Primary = append(manager.Pool.Primary, conn)
	}
	manager.Connections.Store(name, conn)

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
	return newSchema(
		write.Config.Driver,
		&schema.Connection{
			Write: &write.DB,
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
	return newQuery(
		write.Config.Driver,
		&query.Connection{
			Write: &write.DB,
			Read:  &read.DB,
		})
}
