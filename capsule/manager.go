package capsule

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
)

// Add Register a connection with the manager.
func (manager *Manager) Add(name string, driver string, datasource string, readonly bool) (*Manager, error) {

	config := dbal.Config{
		Name:     name,
		Driver:   driver,
		DSN:      datasource,
		ReadOnly: readonly,
	}

	db, err := sqlx.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		DB:     *db,
		Config: &config,
	}

	if config.ReadOnly == true {
		manager.Pool.Readonly = append(manager.Pool.Readonly, conn)
	} else {
		manager.Pool.Primary = append(manager.Pool.Primary, conn)
	}

	manager.Connections.Store(config.Name, conn)
	if Global == nil {
		Global = manager
	}

	return manager, nil
}

// SetAsGlobal Make this connetion instance available globally.
func (manager *Manager) SetAsGlobal() {
	Global = manager
}

// Primary select a primary connection
func (manager *Manager) Primary() (*Connection, error) {
	return manager.Pool.RandPrimary()
}

// ReadOnly select a read-only connection
func (manager *Manager) ReadOnly() (*Connection, error) {
	return manager.Pool.RandReadOnly()
}

// Schema Get a schema builder instance.
func (manager *Manager) Schema() schema.Schema {
	write, err := manager.Primary()
	if err != nil {
		panic(err)
	}

	return schema.Use(&schema.Connection{
		Write:       &write.DB,
		WriteConfig: write.Config,
		Option:      manager.Option,
	})
}

// Query Get a fluent query builder instance.
func (manager *Manager) Query() query.Query {
	write, err := manager.Primary()
	if err != nil {
		panic(err)
	}

	read, err := manager.ReadOnly()
	if err != nil {
		panic(err)
	}

	return query.Use(
		&query.Connection{
			Write:       &write.DB,
			WriteConfig: write.Config,
			Read:        &read.DB,
			ReadConfig:  read.Config,
			Option:      manager.Option,
		})
}

// Close the connections
func (manager *Manager) Close() error {

	messages := []string{}
	manager.Connections.Range(func(key, value any) bool {
		conn, _ := value.(*Connection)
		err := conn.Close()
		if err != nil {
			messages = append(messages, err.Error())
		}
		return true
	})

	if len(messages) > 0 {
		return fmt.Errorf("%s", strings.Join(messages, ";"))
	}
	return nil
}
