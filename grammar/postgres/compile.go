package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
)

// CompileSelect Compile a select query into SQL.
func (grammarSQL Postgres) CompileSelect(query *dbal.Query) string {
	bindingOffset := 0
	return grammarSQL.CompileSelectOffset(query, &bindingOffset)
}

// CompileSelectOffset Compile a select query into SQL.
func (grammarSQL Postgres) CompileSelectOffset(query *dbal.Query, offset *int) string {

	if len(query.Unions) > 0 && query.Aggregate.Func != "" {
		return grammarSQL.CompileUnionAggregate(query)
	}

	sqls := map[string]string{}

	// If the query does not have any columns set, we'll set the columns to the
	// * character to just get all of the columns from the database. Then we
	// can build the query and concatenate all the pieces together as one.
	columns := query.Columns
	if len(columns) == 0 {
		query.AddColumn(grammarSQL.Raw("*"))
	}

	// To compile the query, we'll spin through each component of the query and
	// see if that component exists. If it does we'll just call the compiler
	// function for the component which is responsible for making the SQL.
	sqls["aggregate"] = grammarSQL.CompileAggregate(query, query.Aggregate)
	sqls["columns"] = grammarSQL.CompileColumns(query, query.Columns, offset)
	sqls["from"] = grammarSQL.CompileFrom(query, query.From, offset)
	sqls["joins"] = grammarSQL.CompileJoins(query, query.Joins, offset)
	sqls["wheres"] = grammarSQL.CompileWheres(query, query.Wheres, offset)
	sqls["groups"] = grammarSQL.CompileGroups(query, query.Groups, offset)
	sqls["havings"] = grammarSQL.CompileHavings(query, query.Havings, offset)
	sqls["orders"] = grammarSQL.CompileOrders(query, query.Orders, offset)
	sqls["limit"] = grammarSQL.CompileLimit(query, query.Limit, offset)
	sqls["offset"] = grammarSQL.CompileOffset(query, query.Offset)
	sqls["lock"] = grammarSQL.CompileLock(query, query.Lock)

	sql := ""
	for _, name := range []string{"aggregate", "columns", "from", "joins", "wheres", "groups", "havings", "orders", "limit", "offset", "lock"} {
		segment, has := sqls[name]
		if has && segment != "" {
			sql = sql + segment + " "
		}
	}

	// Compile unions
	if len(query.Unions) > 0 {
		sql = fmt.Sprintf("%s %s", grammarSQL.WrapUnion(sql), grammarSQL.CompileUnions(query, query.Unions, offset))
	}

	// reset columns
	query.Columns = columns
	return strings.Trim(sql, " ")
}

// CompileColumns Compile the "select *" portion of the query.
func (grammarSQL Postgres) CompileColumns(query *dbal.Query, columns []interface{}, bindingOffset *int) string {

	// If the query is actually performing an aggregating select, we will let that
	// compiler handle the building of the select clauses, as it will need some
	// more syntax that is best handled by that function to keep things neat.
	if query.Aggregate.Func != "" {
		return ""
	}

	sql := "select"
	if len(query.DistinctColumns) > 0 {
		sql = fmt.Sprintf("select distinct on (%s)", grammarSQL.Columnize(query.DistinctColumns))
	} else if query.Distinct {
		sql = "select distinct"
	}

	sql = fmt.Sprintf("%s %s", sql, grammarSQL.Columnize(columns))

	for _, col := range columns {
		switch col.(type) {
		case dbal.Select:
			*bindingOffset = *bindingOffset + col.(dbal.Select).Offset
		}
	}

	return sql
}

// CompileWheres Compile an update statement into SQL.
func (grammarSQL Postgres) CompileWheres(query *dbal.Query, wheres []dbal.Where, bindingOffset *int) string {

	// Each type of where clauses has its own compiler function which is responsible
	// for actually creating the where clauses SQL. This helps keep the code nice
	// and maintainable since each clause has a very small method that it uses.
	if len(wheres) == 0 {
		return ""
	}

	clauses := []string{}
	// If we actually have some where clauses, we will strip off the first boolean
	// operator, which is added by the query builders for convenience so we can
	// avoid checking for the first clauses in each of the compilers methods.
	for _, where := range wheres {
		boolen := strings.ToLower(where.Boolean)
		typ := xun.UpperFirst(where.Type)
		// WhereBasic, WhereDate, WhereTime ...
		method := reflect.ValueOf(grammarSQL).MethodByName(fmt.Sprintf("Where%s", typ))
		if method.Kind() == reflect.Func {
			in := []reflect.Value{
				reflect.ValueOf(query),
				reflect.ValueOf(where),
				reflect.ValueOf(bindingOffset),
			}
			out := method.Call(in)
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, out[0].String()))
		}
	}

	conjunction := "where"
	if query.IsJoinClause {
		conjunction = "on"
	}

	return fmt.Sprintf("%s %s", conjunction, grammarSQL.RemoveLeadingBoolean(strings.Join(clauses, " ")))
}

// WhereDate Compile a "where date" clause.
func (grammarSQL Postgres) WhereDate(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}
	return fmt.Sprintf("%s::date %s%s", grammarSQL.Wrap(where.Column), where.Operator, value)
}

// WhereTime Compile a "where time" clause.
func (grammarSQL Postgres) WhereTime(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}
	return fmt.Sprintf("%s::time %s%s", grammarSQL.Wrap(where.Column), where.Operator, value)
}

// WhereDay Compile a "where day" clause.
func (grammarSQL Postgres) WhereDay(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("day", query, where, bindingOffset)
}

// WhereMonth Compile a "where month" clause.
func (grammarSQL Postgres) WhereMonth(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("month", query, where, bindingOffset)
}

// WhereYear Compile a "where year" clause.
func (grammarSQL Postgres) WhereYear(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("year", query, where, bindingOffset)
}

// WhereDateBased  Compile a date based where clause.
func (grammarSQL Postgres) WhereDateBased(typ string, query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}
	return fmt.Sprintf("extract(%s from %s)%s%s", typ, grammarSQL.Wrap(where.Column), where.Operator, value)
}

// CompileLock the lock into SQL.
func (grammarSQL Postgres) CompileLock(query *dbal.Query, lock interface{}) string {
	lockType, ok := lock.(string)
	if ok == false {
		return ""
	} else if lockType == "share" {
		return "for share"
	} else if lockType == "update" {
		return "for update"
	}
	return ""
}
