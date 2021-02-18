package capsule

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
)

// Manager The database manager
type Manager struct {
	Pool        *Pool
	Connections *sync.Map // map[string]*Connection
	Option      *dbal.Option
}

// Pool the connection pool
type Pool struct {
	Primary  []*Connection
	Readonly []*Connection
}

// Connection The database connection
type Connection struct {
	sqlx.DB
	Config *dbal.Config
}
