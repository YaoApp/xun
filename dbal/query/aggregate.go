package query

import (
	"fmt"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Count Retrieve the "count" result of the query.
func (builder *Builder) Count(columns ...interface{}) (int64, error) {
	value, err := builder.numericAggregate("count", columns)
	if err != nil {
		return 0, err
	}
	return value.Int64()
}

// MustCount  Retrieve the "count" result of the query.
func (builder *Builder) MustCount(columns ...interface{}) int64 {
	value, err := builder.Count(columns...)
	utils.PanicIF(err)
	return value
}

// Min Retrieve the minimum value of a given column.
func (builder *Builder) Min(columns ...interface{}) (xun.N, error) {
	return builder.numericAggregate("min", columns)
}

// MustMin Retrieve the minimum value of a given column.
func (builder *Builder) MustMin(columns ...interface{}) xun.N {
	value, err := builder.Min(columns...)
	utils.PanicIF(err)
	return value
}

// Max Retrieve the maximum value of a given column.
func (builder *Builder) Max(columns ...interface{}) (xun.N, error) {
	return builder.numericAggregate("max", columns)
}

// MustMax Retrieve the maximum value of a given column.
func (builder *Builder) MustMax(columns ...interface{}) xun.N {
	value, err := builder.Max(columns...)
	utils.PanicIF(err)
	return value
}

// Sum Retrieve the sum of the values of a given column.
func (builder *Builder) Sum(columns ...interface{}) (xun.N, error) {
	return builder.numericAggregate("sum", columns)
}

// MustSum Retrieve the sum of the values of a given column.
func (builder *Builder) MustSum(columns ...interface{}) xun.N {
	value, err := builder.Sum(columns...)
	utils.PanicIF(err)
	return value
}

// Avg Retrieve the average of the values of a given column.
func (builder *Builder) Avg(columns ...interface{}) (xun.N, error) {
	return builder.numericAggregate("avg", columns)
}

// MustAvg Retrieve the average of the values of a given column.
func (builder *Builder) MustAvg(columns ...interface{}) xun.N {
	value, err := builder.Avg(columns...)
	utils.PanicIF(err)
	return value
}

// Execute an aggregate function on the database.
func (builder *Builder) aggregate(fn string, columns []interface{}) (interface{}, error) {

	qb := builder.clone()
	if len(builder.Query.Unions) == 0 {
		qb.Query.Columns = []interface{}{}
		qb.Query.Bindings["select"] = []interface{}{}
	}
	qb.setAggregate(fn, columns)

	// Debug log
	// fmt.Println("aggregate:", qb.ToSQL())

	rows, err := qb.Get()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("aggregate %s get nothing", fn)
	}

	return rows[0]["aggregate"], nil
}

// numericAggregate Execute a numeric aggregate function on the database.
func (builder *Builder) numericAggregate(fn string, columns []interface{}) (xun.N, error) {
	value, err := builder.aggregate(fn, columns)
	return xun.MakeNum(value), err
}

// setAggregate Set the aggregate property without running the query.
func (builder *Builder) setAggregate(fn string, columns []interface{}) *Builder {
	builder.Query.Aggregate = dbal.Aggregate{
		Func:    fn,
		Columns: columns,
	}
	return builder
}
