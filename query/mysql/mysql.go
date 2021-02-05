package mysql

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/query"
)

// New  Get a fluent query builder instance.
func New(conn *query.Connection) query.Query {
	return &Builder{
		Builder: query.NewBuilder(conn),
	}
}

// Join Add a join clause to the query.
func (builder *Builder) Join() {
	fmt.Printf("\nJoin MySQL: \n===\n%#v\n===\n", builder.Conn)
}
