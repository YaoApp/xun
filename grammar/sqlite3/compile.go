package sqlite3

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileSelect Compile a select query into SQL.
func (grammarSQL SQLite3) CompileSelect(query *dbal.Query) string {
	bindingOffset := 0
	return grammarSQL.CompileSelectOffset(query, &bindingOffset)
}

// CompileSelectOffset Compile a select query into SQL.
func (grammarSQL SQLite3) CompileSelectOffset(query *dbal.Query, offset *int) string {

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
	// sqls["lock"] = grammarSQL.CompileLock()

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
		// fmt.Printf("CompileSelect: %s %s %s %v\n", where.Boolean, where.Type, where.Operator, where.Value)
		boolen := strings.ToLower(where.Boolean)
		switch where.Type {
		case "basic":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereBasic(query, where, bindingOffset)))
			break
		case "date":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereDate(query, where, bindingOffset)))
			break
		case "time":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereTime(query, where, bindingOffset)))
			break
		case "day":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereDay(query, where, bindingOffset)))
			break
		case "month":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereMonth(query, where, bindingOffset)))
			break
		case "year":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereYear(query, where, bindingOffset)))
			break
		case "raw":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereRaw(query, where, bindingOffset)))
			break
		case "null":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereNull(query, where, bindingOffset)))
			break
		case "notnull":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereNotNull(query, where, bindingOffset)))
			break
		case "between":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereBetween(query, where, bindingOffset)))
			break
		case "in":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereIn(query, where, bindingOffset)))
			break
		case "column":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereColumn(query, where, bindingOffset)))
			break
		case "sub":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereSub(query, where, bindingOffset)))
			break
		case "exists":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.whereExists(query, where, bindingOffset)))
			break
		case "nested":
			clauses = append(clauses, fmt.Sprintf("%s %s", boolen, grammarSQL.WhereNested(query, where, bindingOffset)))
			break
		}
	}

	conjunction := "where"
	if query.IsJoinClause {
		conjunction = "on"
	}

	return fmt.Sprintf("%s %s", conjunction, grammarSQL.RemoveLeadingBoolean(strings.Join(clauses, " ")))
}

// Compile a where (not) exists clause.
func (grammarSQL SQLite3) whereExists(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	exists := "exists"
	if where.Not {
		exists = "not exists"
	}
	selectSQL := grammarSQL.CompileSelectOffset(where.Query, bindingOffset)
	return fmt.Sprintf("%s (%s)", exists, selectSQL)
}

// WhereDate Compile a "where date" clause.
func (grammarSQL SQLite3) WhereDate(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("%Y-%m-%d", query, where, bindingOffset)
}

// WhereTime Compile a "where time" clause.
func (grammarSQL SQLite3) WhereTime(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("time", query, where, bindingOffset)
}

// WhereTime Compile a "where day" clause.
func (grammarSQL SQLite3) whereDay(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("day", query, where, bindingOffset)
}

// whereMonth Compile a "where month" clause.
func (grammarSQL SQLite3) whereMonth(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("month", query, where, bindingOffset)
}

// whereYear Compile a "where year" clause.
func (grammarSQL SQLite3) whereYear(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("year", query, where, bindingOffset)
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

	return fmt.Sprintf("strftime('%s',%s) %s cast(%s as text)", typ, grammarSQL.Wrap(where.Column), where.Operator, value)
}
