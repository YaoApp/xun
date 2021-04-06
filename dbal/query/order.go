package query

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// OrderBy Add an "order by" clause to the query.
func (builder *Builder) OrderBy(column interface{}, args ...string) Query {
	offset := 0
	direction := "asc"
	orderName := "order"

	if len(args) > 0 {
		direction = strings.ToLower(args[0])
	}
	if len(builder.Query.Unions) > 0 {
		orderName = "unionOrder"
	}

	if !utils.StringHave([]string{"asc", "desc"}, direction) {
		panic(fmt.Errorf(`Order direction must be "asc" or "desc`))
	}

	if builder.isQueryable(column) {
		sql := ""
		bindings := []interface{}{}
		sql, bindings, offset = builder.createSub(column)
		column = dbal.Raw(fmt.Sprintf("(%s)", sql))
		builder.Query.AddBinding(orderName, bindings)
	}

	order := dbal.Order{
		Type:      "basic",
		Column:    column,
		Direction: direction,
		Offset:    offset,
	}

	if orderName == "unionOrder" {
		builder.Query.UnionOrders = append(builder.Query.UnionOrders, order)
	} else if orderName == "order" {
		builder.Query.Orders = append(builder.Query.Orders, order)
	}
	return builder
}

// OrderByDesc Add a descending "order by xxx desc" clause to the query.
func (builder *Builder) OrderByDesc(column interface{}) Query {
	return builder.OrderBy(column, "desc")
}

// OrderByRaw Add a raw "order by" clause to the query.
func (builder *Builder) OrderByRaw(sql string, bindings ...interface{}) Query {
	order := dbal.Order{
		Type: "raw",
		SQL:  sql,
	}
	if len(builder.Query.Unions) > 0 {
		builder.Query.UnionOrders = append(builder.Query.UnionOrders, order)
		if len(bindings) > 0 {
			builder.Query.AddBinding("unionOrder", bindings)
		}
	} else {
		builder.Query.Orders = append(builder.Query.Orders, order)
		if len(bindings) > 0 {
			builder.Query.AddBinding("order", bindings)
		}
	}
	return builder
}

// Latest Add an "order by" clause for a timestamp to the query. @todo
func (builder *Builder) Latest() {
}

// Oldest Add an "order by" clause for a timestamp to the query. @todo
func (builder *Builder) Oldest() {
}

// InRandomOrder Put the query's results in random order. @todo
func (builder *Builder) InRandomOrder() {
}

// Reorder Remove all existing orders and optionally add a new order. @todo
func (builder *Builder) Reorder() {
}
