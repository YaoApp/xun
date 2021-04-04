package query

import "github.com/yaoapp/xun/dbal"

// GroupBy Add a "group by" clause to the query.
func (builder *Builder) GroupBy(groups ...interface{}) Query {
	builder.Query.Groups = append(builder.Query.Groups, groups...)
	return builder
}

// GroupByRaw Add a raw groupBy clause to the query.
func (builder *Builder) GroupByRaw(expression string, bindings ...interface{}) Query {
	builder.Query.Groups = []interface{}{}
	builder.Query.Groups = append(builder.Query.Groups, dbal.Raw(expression))
	if len(bindings) > 0 {
		builder.Query.AddBinding("groupBy", bindings)
	}
	return builder
}

// Having Add a "having" clause to the query.
func (builder *Builder) Having() {
}

// OrHaving Add an "or having" clause to the query.
func (builder *Builder) OrHaving() {
}

// HavingBetween Add a "having between " clause to the query.
func (builder *Builder) HavingBetween() {
}

// OrHavingBetween Add an "or having between" clause to the query.
func (builder *Builder) OrHavingBetween() {
}

// HavingRaw Add a raw having clause to the query.
func (builder *Builder) HavingRaw() {
}

// OrHavingRaw Add a raw or having clause to the query.
func (builder *Builder) OrHavingRaw() {
}
