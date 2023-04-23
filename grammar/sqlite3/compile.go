package sqlite3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
)

// CompileSelect Compile a select query into SQL.
func (grammarSQL SQLite3) CompileSelect(query *dbal.Query) string {
	bindingOffset := 0
	return grammarSQL.CompileSelectOffset(query, &bindingOffset)
}

// CompileSelectOffset Compile a select query into SQL.
func (grammarSQL SQLite3) CompileSelectOffset(query *dbal.Query, offset *int) string {

	// SQL STMT
	if query.SQL != "" {
		if !strings.Contains(query.SQL, "limit") && !strings.Contains(query.SQL, "offset") {
			limit := grammarSQL.CompileLimit(query, query.Limit, offset)
			offset := grammarSQL.CompileOffset(query, query.Offset)
			return strings.TrimSpace(fmt.Sprintf("%s %s %s", query.SQL, limit, offset))
		}
		return query.SQL
	}

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

// CompileWheres Compile an update statement into SQL.
func (grammarSQL SQLite3) CompileWheres(query *dbal.Query, wheres []dbal.Where, bindingOffset *int) string {

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
func (grammarSQL SQLite3) WhereDate(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%Y-%m-%d", query, where, bindingOffset)
}

// WhereTime Compile a "where time" clause.
func (grammarSQL SQLite3) WhereTime(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%H:%M:%S", query, where, bindingOffset)
}

// WhereDay Compile a "where day" clause.
func (grammarSQL SQLite3) WhereDay(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%d", query, where, bindingOffset)
}

// WhereMonth Compile a "where month" clause.
func (grammarSQL SQLite3) WhereMonth(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%m", query, where, bindingOffset)
}

// WhereYear Compile a "where year" clause.
func (grammarSQL SQLite3) WhereYear(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%Y", query, where, bindingOffset)
}

// WhereDateBased  Compile a date based where clause.
func (grammarSQL SQLite3) WhereDateBased(typ string, query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}

	return fmt.Sprintf("strftime('%s',%s) %s cast(%s as text)", typ, grammarSQL.Wrap(where.Column, false), where.Operator, value)
}

// CompileLock the lock into SQL.
func (grammarSQL SQLite3) CompileLock(query *dbal.Query, lock interface{}) string {
	return ""
}
