package query

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// From set the table which the query is targeting.
func (builder *Builder) From(from string) Query {
	name := dbal.NewName(from, builder.Conn.Option.Prefix)
	builder.Query.From = dbal.From{
		Type:   "basic",
		Alias:  name.Alias,
		Name:   name,
		Offset: 0,
	}
	return builder
}

// FromRaw Add a raw from clause to the query.
func (builder *Builder) FromRaw(sql string, bindings ...interface{}) Query {
	builder.Query.From = dbal.From{
		Type:   "raw",
		SQL:    sql,
		Offset: len(bindings),
	}
	builder.Query.AddBinding("from", bindings)
	return builder
}

// FromSub Makes "from" fetch from a subquery.
func (builder *Builder) FromSub(qb interface{}, as string) Query {
	sub, bindings, fromOffset := builder.createSub(qb)
	segment := builder.parseSub(sub)
	builder.Query.From = dbal.From{
		Type:   "sub",
		Alias:  as,
		SQL:    fmt.Sprintf("(%s)", segment),
		Offset: fromOffset,
	}
	builder.Query.AddBinding("from", bindings)
	return builder
}
