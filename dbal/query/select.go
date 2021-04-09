package query

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// Select Set the columns to be selected.
// Select("field1", "field2")
// Select("field1", "field2 as f2")
// Select("field1", dbal.Raw("Count(id) as v"))
func (builder *Builder) Select(columns ...interface{}) Query {
	builder.Query.Columns = []interface{}{}
	builder.Query.Bindings["select"] = []interface{}{}
	if len(columns) == 0 {
		builder.Query.AddColumn(dbal.Raw("*"))
	}

	columns = builder.prepareColumns(columns...)
	for _, column := range columns {
		builder.addSelect(column)
	}
	return builder
}

// SelectRaw Add a new "raw" select expression to the query.
func (builder *Builder) SelectRaw(expression string, bindings ...interface{}) Query {
	builder.addSelect(dbal.Raw(expression))
	builder.Query.AddBinding("select", bindings)
	return builder
}

// SelectSub Add a subselect expression to the query.
func (builder *Builder) SelectSub(qb interface{}, as string) Query {
	sub, bindings, selectOffset := builder.createSub(qb)
	sql := builder.parseSub(sub)
	column := dbal.Select{
		Type:   "sub",
		Alias:  as,
		SQL:    fmt.Sprintf("(%s)", sql),
		Offset: selectOffset - 1,
	}
	builder.addSelect(column)
	builder.Query.AddBinding("select", bindings)
	return builder
}

// Distinct Force the query to only return distinct results.
func (builder *Builder) Distinct(args ...interface{}) Query {
	if len(args) > 0 {
		if _, ok := args[0].(bool); ok {
			builder.Query.Distinct = args[0].(bool)
		} else if _, ok := args[0].([]interface{}); ok {
			builder.Query.Distinct = true
			builder.Query.DistinctColumns = args[0].([]interface{})
		} else if _, ok := args[0].([]string); ok {
			builder.Query.Distinct = true
			builder.Query.DistinctColumns = []interface{}{}
			for _, arg := range args[0].([]string) {
				builder.Query.DistinctColumns = append(builder.Query.DistinctColumns, arg)
			}
		} else {
			builder.Query.Distinct = true
			builder.Query.DistinctColumns = args
		}
	} else {
		builder.Query.Distinct = true
	}
	return builder
}

// addSelect Add a new select column to the query.
func (builder *Builder) addSelect(column interface{}) {
	switch column.(type) {
	case string:
		if strings.Contains(column.(string), ",") {
			columns := strings.Split(column.(string), ",")
			for _, col := range columns {
				col = strings.Trim(col, " ")
				builder.addSelect(col)
			}
		} else {
			builder.Query.Columns = append(builder.Query.Columns, column)
		}
	default:
		builder.Query.Columns = append(builder.Query.Columns, column)
	}
}
