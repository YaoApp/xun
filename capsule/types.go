package capsule

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal/schema"
)

// Manager The database manager
type Manager struct {
	Pool        *Pool
	Connections *sync.Map // map[string]*Connection
	Config      *schema.Config
}

// Pool the connection pool
type Pool struct {
	Primary  []*Connection
	Readonly []*Connection
}

// Connection The database connection
type Connection struct {
	sqlx.DB
	Config *schema.ConnConfig
}
