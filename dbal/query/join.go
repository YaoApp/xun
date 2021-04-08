package query

import (
	"reflect"

	"github.com/yaoapp/xun/dbal"
)

// Join  Add a join clause to the query.
func (builder *Builder) Join(table string, first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	return builder.join(dbal.NewName(table), "", first, operator, second, "inner", "on", 0)
}

// JoinSub Add a subquery join clause to the query.
func (builder *Builder) JoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	return builder.joinSub(qb, alias, first, operator, second, "inner", "on", 0)
}

// LeftJoin Add a left join to the query.
func (builder *Builder) LeftJoin(table string, first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	return builder.join(dbal.NewName(table), "", first, operator, second, "left", "on", 0)
}

// LeftJoinSub Add a subquery left join to the query.
func (builder *Builder) LeftJoinSub(qb interface{}, alias string, first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	return builder.joinSub(qb, alias, first, operator, second, "left", "on", 0)
}

// RightJoin Add a right join to the query.
func (builder *Builder) RightJoin() {
}

// CrossJoin Add a "cross join" clause to the query.
func (builder *Builder) CrossJoin() {
}

// On Add an "on" clause to the join.
func (builder *Builder) On(first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	join := builder.joinOn(first, operator, second, "and", 0)
	builder.Query.Joins = append(builder.Query.Joins, join)
	return builder
}

// OrOn Add an "or on" clause to the join.
func (builder *Builder) OrOn(first interface{}, args ...interface{}) Query {
	operator, second := builder.joinPrepare(args...)
	join := builder.joinOn(first, operator, second, "or", 0)
	builder.Query.Joins = append(builder.Query.Joins, join)
	return builder
}

// forJoinClause Create a new query instance for a join clause.
func (builder *Builder) forJoinClause() *Builder {
	new := builder.new()
	new.Query.IsJoinClause = true
	return new
}

func (builder *Builder) joinPrepare(args ...interface{}) (string, interface{}) {
	var operator = "="
	var second interface{} = nil

	if len(args) == 1 {
		second = args[0]
	}

	if len(args) > 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
		operator = args[0].(string)
	}

	if len(args) > 1 {
		second = args[1]
	}

	return operator, second
}

//JoinWhere Add a "join where" clause to the query.
func (builder *Builder) JoinWhere() {
}

// LeftJoinWhere Add a "join where" clause to the query.
func (builder *Builder) LeftJoinWhere() {
}

// RightJoinWhere Add a "right join where" clause to the query.
func (builder *Builder) RightJoinWhere() {
}

//Join Add a join clause to the query.
func (builder *Builder) join(table interface{}, alias string, first interface{}, operator string, second interface{}, typ string, method string, offset int) Query {
	if method == "on" {
		qb := builder.forJoinClause()
		join := qb.joinOn(first, operator, second, "and", offset)
		join.Type = typ
		join.Name = table
		if alias != "" {
			join.Alias = alias
			join.SQL = dbal.Raw(table).GetValue()
		}
		builder.Query.Joins = append(builder.Query.Joins, join)
		builder.Query.AddBinding("join", qb.GetBindings())
	}
	return builder
}

// joinSub Add a subquery join clause to the query.
func (builder *Builder) joinSub(qb interface{}, alias string, first interface{}, operator string, second interface{}, typ string, method string, offset int) Query {
	sub, bindings, joinOffset := builder.createSub(qb)
	builder.Query.AddBinding("join", bindings)
	return builder.join(sub, alias, first, operator, second, typ, method, joinOffset+offset)
}

func (builder *Builder) joinOn(first interface{}, operator string, second interface{}, boolean string, offset int) dbal.Join {

	if builder.isClosure(first) {
		builder.whereNested(first.(func(qb Query)), boolean)
	} else {
		builder.whereColumn(first, operator, second, boolean, offset)
	}

	return dbal.Join{
		Query:  builder.Query,
		Offset: offset,
	}
}
