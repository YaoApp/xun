package query

import (
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
	for _, column := range columns {
		builder.Query.Columns = append(builder.Query.Columns, column)
	}
	return builder
}

// SelectRaw Add a new "raw" select expression to the query.
func (builder *Builder) SelectRaw(expression string, bindings ...interface{}) Query {
	builder.addSelect(dbal.Raw(expression))
	if len(bindings) > 0 {
		builder.Query.AddBinding("select", bindings)
	}
	return builder
}

// addSelect Add a new select column to the query.
func (builder *Builder) addSelect(column interface{}) {
	builder.Query.Columns = append(builder.Query.Columns, column)
}

// SelectSub Add a subselect expression to the query.
func (builder *Builder) SelectSub() {
}

// Distinct Force the query to only return distinct results.
func (builder *Builder) Distinct() {
}
