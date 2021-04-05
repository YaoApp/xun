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
func (builder *Builder) Having(column interface{}, args ...interface{}) Query {

	typ := "basic"

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, value, boolean, offset := builder.prepareArgs(args...)

	builder.Query.Havings = append(builder.Query.Havings, dbal.Having{
		Type:     typ,
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
		Offset:   offset,
		Not:      false,
	})

	if !builder.isExpression(value) {
		builder.Query.AddBinding("having", builder.flattenValue(value))
	}

	return builder
}

// OrHaving Add an "or having" clause to the query.
func (builder *Builder) OrHaving(column interface{}, args ...interface{}) Query {
	operator, value, _, offset := builder.prepareArgs(args...)
	return builder.Having(column, operator, value, "or", offset)
}

// HavingBetween Add a "having between " clause to the query.
func (builder *Builder) HavingBetween(column interface{}, values []interface{}, args ...interface{}) Query {

	boolean := "and"
	not := false
	offset := 1
	if len(args) > 0 {
		if _, ok := args[0].(string); ok {
			boolean = args[0].(string)
		}
	}

	if len(args) > 1 {
		if _, ok := args[1].(bool); ok {
			not = args[1].(bool)
		}
	}

	if len(args) > 2 {
		if _, ok := args[2].(int); ok {
			offset = args[2].(int)
		}
	}

	builder.Query.Havings = append(builder.Query.Havings, dbal.Having{
		Type:    "between",
		Column:  column,
		Not:     not,
		Boolean: boolean,
		Values:  values,
		Offset:  offset,
	})

	values = builder.cleanBindings(values)
	if len(values) > 2 {
		values = values[0:2]
	}

	if len(values) > 0 {
		builder.Query.AddBinding("having", values)
	}
	return builder
}

// OrHavingBetween Add a "having between " clause to the query.
func (builder *Builder) OrHavingBetween(column interface{}, values []interface{}, args ...interface{}) Query {
	not := false
	offset := 1
	if len(args) > 0 {
		if _, ok := args[0].(bool); ok {
			not = args[0].(bool)
		}
	}
	if len(args) > 1 {
		if _, ok := args[1].(int); ok {
			offset = args[1].(int)
		}
	}
	return builder.HavingBetween(column, values, "or", not, offset)
}

// HavingRaw Add a raw having clause to the query.
func (builder *Builder) HavingRaw(sql string, bindings ...interface{}) Query {
	builder.Query.Havings = append(builder.Query.Havings, dbal.Having{
		Type:    "raw",
		SQL:     sql,
		Boolean: "and",
		Offset:  0,
	})
	if len(bindings) > 0 {
		builder.Query.AddBinding("having", bindings)
	}
	return builder
}

// OrHavingRaw Add a raw or having clause to the query.
func (builder *Builder) OrHavingRaw(sql string, bindings ...interface{}) Query {
	builder.Query.Havings = append(builder.Query.Havings, dbal.Having{
		Type:    "raw",
		SQL:     sql,
		Boolean: "or",
		Offset:  0,
	})
	if len(bindings) > 0 {
		builder.Query.AddBinding("having", bindings)
	}
	return builder
}
