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
		sub, bindings, subQueryOffset := builder.createSub(column)
		sql := builder.parseSub(sub)
		offset = subQueryOffset
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
func (builder *Builder) Reorder(args ...interface{}) Query {
	builder.Query.Orders = []dbal.Order{}
	builder.Query.UnionOrders = []dbal.Order{}
	builder.Query.Bindings["order"] = []interface{}{}
	builder.Query.Bindings["unionOrder"] = []interface{}{}

	if len(args) == 1 {
		return builder.OrderBy(args[0])
	} else if len(args) == 2 {
		direction := "asc"
		if _, ok := args[0].(string); ok {
			direction = args[0].(string)
		}
		return builder.OrderBy(args[0], direction)
	}

	return builder
}

// removeExistingOrdersFor Get an array with all orders with a given column removed.
func (builder *Builder) removeExistingOrdersFor(column interface{}) []dbal.Order {
	index := 0
	remove := false
	orders := builder.Query.Orders
	for i, order := range orders {
		if order.Column == column {
			remove = true
			index = i
			break
		}
	}
	if remove {
		orders = append(orders[:index], orders[index+1:]...)
	}
	return orders
}

// Throw an exception if the query doesn't have an orderBy clause.
func (builder *Builder) enforceOrderBy() {
	if len(builder.Query.Orders) == 0 && len(builder.Query.UnionOrders) == 0 {
		panic("You must specify an orderBy clause when using this function.")
	}
}
