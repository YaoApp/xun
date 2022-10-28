package capsule

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Load mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
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

// Add Register a connection with the manager.
func Add(name string, driver string, dsn string) (*Manager, error) {
	return New().Add(name, driver, dsn, false)
}

// AddRead Register a readonly connection with the manager.
func AddRead(name string, driver string, dsn string) (*Manager, error) {
	return New().Add(name, driver, dsn, true)
}

// Schema Get a schema builder instance.
func Schema() schema.Schema {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Schema()
}

// Query Get a fluent query builder instance.
func Query() query.Query {
	if Global == nil {
		err := errors.New("the global capsule not set")
		panic(err)
	}
	return Global.Query()
}

// ************************************************************
// THE FOLLOWING LINES WILL BE DEPRECATED
// ************************************************************

// NewWithOption Create a database manager instance using the given option.
func NewWithOption(option dbal.Option) *Manager {
	manager := New()
	manager.SetOption(option)
	return manager
}

// AddConn Register a connection with the manager.
func AddConn(name string, driver string, datasource string, timeout ...time.Duration) *Manager {
	return New().AddConn(name, driver, datasource, timeout...)
}

// AddConn Register a connection with the manager.
func (manager *Manager) AddConn(name string, driver string, datasource string, timeout ...time.Duration) *Manager {
	manager.AddConnection(name, driver, datasource, false, timeout...)
	return manager
}

// AddReadConn Register a readonly connection with the manager.
func AddReadConn(name string, driver string, datasource string, timeout ...time.Duration) *Manager {
	return New().AddReadConn(name, driver, datasource, timeout...)
}

// AddReadConn Register a readonly with the manager.
func (manager *Manager) AddReadConn(name string, driver string, datasource string, timeout ...time.Duration) *Manager {
	manager.AddConnection(name, driver, datasource, true, timeout...)
	return manager
}

// SetOption set the database manager as the given value
func (manager *Manager) SetOption(option dbal.Option) {
	manager.Option = &option
}

// AddConnection Register a connection with the manager.
func (manager *Manager) AddConnection(name string, driver string, datasource string, readonly bool, timeouts ...time.Duration) *Manager {
	config := dbal.Config{
		Name:     name,
		Driver:   driver,
		DSN:      datasource,
		ReadOnly: readonly,
	}

	db := sqlx.MustOpen(config.Driver, config.DSN)

	// Cheking database connection
	timeout := 1 * time.Second
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		err := db.PingContext(ctx)
		if err != nil {
			panic(fmt.Sprintf("Connection timeout %s (%s: %s)", timeout, config.Driver, config.DSN))
		}
		cancel()
	}()

	<-ctx.Done()
	conn := &Connection{
		DB:     *db,
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
