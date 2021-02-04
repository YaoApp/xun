package sqlserver

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/query"
)

// Table  Get a fluent query builder instance.
func Table() query.Query {
	return &Builder{
		Builder: query.NewBuilder(),
	}
}

// Join Add a join clause to the query.
func (builder *Builder) Join() { fmt.Printf("SQLServer JOIN\n") }
