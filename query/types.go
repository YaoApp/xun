package query

import "github.com/yaoapp/xun/dbal/query"

// Builder a query builder
type Builder struct{ query.Builder }

// Query The database Query interface
type Query interface{ query.Query }
