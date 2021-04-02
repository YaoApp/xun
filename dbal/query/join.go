package query

import (
	"reflect"

	"github.com/yaoapp/xun/dbal"
)

//Join  Add a join clause to the query.
func (builder *Builder) Join(table string, first interface{}, args ...interface{}) Query {
	operator, second, typ, method := builder.joinPrepare(args...)
	if method == "on" {
		qb := builder.forJoinClause()
		join := qb.joinOn(table, first, operator, second, typ, "and")
		builder.Query.Joins = append(builder.Query.Joins, join)
		builder.Query.AddBinding("join", qb.GetBindings())
	}
	return builder
}

// forJoinClause Create a new query instance for a join clause.
func (builder *Builder) forJoinClause() *Builder {
	new := builder.new()
	new.Query.IsJoinClause = true
	return new
}

func (builder *Builder) joinOn(table string, first interface{}, operator string, second interface{}, typ string, boolean string) dbal.Join {

	if builder.isClosure(first) {
		builder.whereNested(first.(func(qb Query)), boolean)
	} else {
		builder.WhereColumn(first, operator, second, boolean)
	}

	return dbal.Join{
		Type:  typ,
		Table: dbal.NewName(table),
		Query: builder.Query,
	}
}

func (builder *Builder) joinPrepare(args ...interface{}) (string, interface{}, string, string) {
	var operator = ""
	var second interface{} = nil
	var typ = "inner"
	var method = "on"

	if len(args) > 0 && reflect.TypeOf(args[0]).Kind() == reflect.String {
		operator = args[0].(string)
	}

	if len(args) > 1 {
		second = args[1]
	}

	if len(args) > 2 && reflect.TypeOf(args[2]).Kind() == reflect.String {
		typ = args[2].(string)
	}

	if len(args) > 3 && reflect.TypeOf(args[3]).Kind() == reflect.String {
		typ = args[3].(string)
	}

	return operator, second, typ, method
}

//JoinWhere Add a "join where" clause to the query.
func (builder *Builder) JoinWhere() {
}

// LeftJoin Add a left join to the query.
func (builder *Builder) LeftJoin() {
}

// LeftJoinWhere Add a "join where" clause to the query.
func (builder *Builder) LeftJoinWhere() {
}

// RightJoin Add a right join to the query.
func (builder *Builder) RightJoin() {
}

// RightJoinWhere Add a "right join where" clause to the query.
func (builder *Builder) RightJoinWhere() {
}

// CrossJoin Add a "cross join" clause to the query.
func (builder *Builder) CrossJoin() {
}
