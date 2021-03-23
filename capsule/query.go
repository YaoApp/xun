package capsule

import (
	"github.com/yaoapp/xun/dbal/query"
)

// newQuery Get a query builder instance.
func newQuery(driver string, conn *query.Connection) query.Query {
	return query.Use(conn)
}
