package capsule

import (
	"errors"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/query"
	"github.com/yaoapp/xun/schema"
)

// New Create a database manager instance.
func New() *Manager {
	return &Manager{
		Pools:       &sync.Map{},
		Connections: &sync.Map{},
		Current:     nil,
	}
}

// AddConnection Register a connection with the manager.
func (manager *Manager) AddConnection(config Config) *Manager {

	name := "main"
	if config.Name != "" {
		name = config.Name
	}

	group := "default"
	if config.Group != "" {
		group = config.Group
	}

	conn := &Connection{
		DB: *sqlx.MustOpen(config.DriverName(), config.DataSource()),
	}

	p, has := manager.Pools.Load(group)
	var pool Pool
	if !has {
		pool = Pool{
			Primary:  []*Connection{},
			Readonly: []*Connection{},
		}
	} else {
		pool = p.(Pool)
	}

	if config.ReadOnly == true {
		pool.Readonly = append(pool.Readonly, conn)
	} else {
		pool.Primary = append(pool.Primary, conn)
	}
	manager.Connections.Store(group+"."+name, conn)

	if manager.Current == nil {
		manager.Current = conn
	}

	return manager
}

// GetConnection Get a registered connection instance.
func (manager *Manager) GetConnection(name string, group ...string) *Connection {
	groupName := "default"
	if len(group) > 0 {
		groupName = group[0]
	}
	c, has := manager.Connections.Load(groupName + "." + name)
	conn := c.(*Connection)
	if !has {
		err := errors.New("the connection " + groupName + "." + name + " is not registered")
		panic(err)
	}
	return conn
}

// SetAsGlobal Make this connetion instance available globally.
func (manager *Manager) SetAsGlobal(name string, group ...string) {
	conn := manager.GetConnection(name, group...)
	manager.Current = conn
}

// Get Get current connection
func (manager *Manager) Get() *Connection {
	return manager.Current
}

// Schema Get a schema builder instance.
func (manager *Manager) Schema() schema.Schema {
	return schema.New()
}

// Table Get a fluent query builder instance.
func (manager *Manager) Table() query.Query {
	return query.Table()
}
