package sql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
)

// CompileSelect Compile a select query into SQL.
func (grammarSQL SQL) CompileSelect(query *dbal.Query) string {
	bindingOffset := 0
	return grammarSQL.CompileSelectOffset(query, &bindingOffset)
}

// CompileSelectOffset Compile a select query into SQL.
func (grammarSQL SQL) CompileSelectOffset(query *dbal.Query, offset *int) string {

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

// CompileUnionAggregate Compile a union aggregate query into SQL.
func (grammarSQL SQL) CompileUnionAggregate(query *dbal.Query) string {
	qb := &(*query)
	sql := grammarSQL.CompileAggregate(query, query.Aggregate)
	qb.Aggregate = dbal.Aggregate{}
	return fmt.Sprintf("%s from (%s) as %s", sql, grammarSQL.CompileSelect(qb), grammarSQL.WrapTable("temp_table"))
}

// CompileAggregate Compile an aggregated select clause.
func (grammarSQL SQL) CompileAggregate(query *dbal.Query, aggregate dbal.Aggregate) string {

	if aggregate.Func == "" {
		return ""
	}

	column := grammarSQL.Columnize(aggregate.Columns)

	// If the query has a "distinct" constraint and we're not asking for all columns
	// we need to prepend "distinct" onto the column name so that the query takes
	// it into account when it performs the aggregating operations on the data.
	if len(query.DistinctColumns) > 0 {
		column = fmt.Sprintf("distinct %s", grammarSQL.Columnize(query.DistinctColumns))
	} else if query.Distinct && column != "*" {
		column = fmt.Sprintf("distinct %s", column)
	}
	return fmt.Sprintf("select %s(%s) as aggregate", aggregate.Func, column)
}

// CompileUnions Compile the "union" queries attached to the main query.
func (grammarSQL SQL) CompileUnions(query *dbal.Query, unions []dbal.Union, offset *int) string {
	sql := ""
	for _, union := range unions {
		sql = sql + grammarSQL.CompileUnion(query, union, offset)
	}

	// unionOrders
	if len(query.UnionOrders) > 0 {
		sql = fmt.Sprintf("%s %s", sql, grammarSQL.CompileOrders(query, query.UnionOrders, offset))
	}

	// unionLimit
	if query.UnionLimit >= 0 {
		sql = fmt.Sprintf("%s %s", sql, grammarSQL.CompileLimit(query, query.UnionLimit, offset))
	}

	// unionOffset
	if query.UnionOffset >= 0 {
		sql = fmt.Sprintf("%s %s", sql, grammarSQL.CompileOffset(query, query.UnionOffset))
	}

	return strings.TrimPrefix(sql, " ")
}

// CompileUnion Compile a single union statement.
func (grammarSQL SQL) CompileUnion(query *dbal.Query, union dbal.Union, offset *int) string {
	conjunction := "union "
	if union.All {
		conjunction = "union all "
	}
	return fmt.Sprintf("%s%s", conjunction, grammarSQL.WrapUnion(grammarSQL.CompileSelectOffset(union.Query, offset)))
}

// CompileJoins Compile the "join" portions of the query.
func (grammarSQL SQL) CompileJoins(query *dbal.Query, joins []dbal.Join, offset *int) string {
	sql := ""
	for _, join := range joins {
		table := grammarSQL.WrapTable(join.Name)
		if join.SQL != nil && join.Alias != "" {
			sql := grammarSQL.CompileSub(join.SQL, offset)
			table = fmt.Sprintf("(%s) as %s", sql, join.Alias)
		}
		nestedJoins := " "
		if len(join.Query.Joins) > 0 {
			nestedJoins = grammarSQL.CompileJoins(query, join.Query.Joins, offset)
		}
		tableAndNestedJoins := table
		if len(join.Query.Joins) > 0 {
			tableAndNestedJoins = fmt.Sprintf("(%s%s)", table, nestedJoins)
		}

		return strings.Trim(
			fmt.Sprintf("%s join %s %s", join.Type, tableAndNestedJoins, grammarSQL.CompileWheres(join.Query, join.Query.Wheres, offset)),
			" ",
		)
	}
	return sql
}

// CompileSub Parse the subquery into SQL and bindings.
func (grammarSQL SQL) CompileSub(sub interface{}, offset *int) string {
	switch sub.(type) {
	case *dbal.Query:
		query := sub.(*dbal.Query)
		return grammarSQL.CompileSelectOffset(query, offset)
	case dbal.Expression:
		return sub.(dbal.Expression).GetValue()
	case string:
		return sub.(string)
	}
	panic(fmt.Errorf("a subquery must be a query builder instance, a Closure, or a string"))
}

// CompileColumns Compile the "select *" portion of the query.
func (grammarSQL SQL) CompileColumns(query *dbal.Query, columns []interface{}, bindingOffset *int) string {

	// If the query is actually performing an aggregating select, we will let that
	// compiler handle the building of the select clauses, as it will need some
	// more syntax that is best handled by that function to keep things neat.
	if query.Aggregate.Func != "" {
		return ""
	}

	sql := "select"
	if query.Distinct {
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

//CompileFrom  Compile the "from" portion of the query.
func (grammarSQL SQL) CompileFrom(query *dbal.Query, from dbal.From, bindingOffset *int) string {
	if from.Type == "" {
		return ""
	}

	sql := ""
	if from.Type == "raw" {
		sql = fmt.Sprintf("from %s", from.SQL)
	} else if from.Type == "sub" {
		if from.Alias != "" {
			sql = fmt.Sprintf("from %s as %s", from.SQL, grammarSQL.ID(from.Alias))
		} else {
			sql = fmt.Sprintf("from %s", from.SQL)
		}
	} else {
		sql = fmt.Sprintf("from %s", grammarSQL.WrapTable(from))
	}
	*bindingOffset = *bindingOffset + from.Offset
	return sql
}

// CompileGroups Compile the "group by" portions of the query.
func (grammarSQL SQL) CompileGroups(query *dbal.Query, groups []interface{}, bindingOffset *int) string {
	if len(groups) == 0 {
		return ""
	}
	return fmt.Sprintf("group by %s", grammarSQL.Columnize(groups))
}

// CompileHavings Compile the "having" portions of the query.
func (grammarSQL SQL) CompileHavings(query *dbal.Query, havings []dbal.Having, bindingOffset *int) string {
	clauses := []string{}
	for _, having := range havings {
		clauses = append(clauses, grammarSQL.CompileHaving(query, having, bindingOffset))
	}
	if len(clauses) == 0 {
		return ""
	}
	return fmt.Sprintf("having %s", grammarSQL.RemoveLeadingBoolean(strings.Join(clauses, " ")))
}

// CompileHaving Compile a single having clause.
func (grammarSQL SQL) CompileHaving(query *dbal.Query, having dbal.Having, bindingOffset *int) string {
	// If the having clause is "raw", we can just return the clause straight away
	// without doing any more processing on it. Otherwise, we will compile the
	// clause into SQL based on the components that make it up from builder.
	if having.Type == "raw" {
		return fmt.Sprintf("%s %s", having.Boolean, having.SQL)
	} else if having.Type == "between" {
		return grammarSQL.HavingBetween(query, having, bindingOffset)
	}
	return grammarSQL.HavingBasic(query, having, bindingOffset)
}

// HavingBasic  Compile a basic having clause.
func (grammarSQL SQL) HavingBasic(query *dbal.Query, having dbal.Having, bindingOffset *int) string {
	if !dbal.IsExpression(having.Value) {
		*bindingOffset = *bindingOffset + having.Offset
	}
	column := grammarSQL.Wrap(having.Column)
	parameter := grammarSQL.Parameter(having.Value, *bindingOffset)

	return fmt.Sprintf("%s %s %s %s", having.Boolean, column, having.Operator, parameter)
}

// HavingBetween Compile a "between" having clause.
func (grammarSQL SQL) HavingBetween(query *dbal.Query, having dbal.Having, bindingOffset *int) string {
	if len(having.Values) != 2 {
		panic(fmt.Errorf("The given values must have two items"))
	}
	between := "between"
	if having.Not {
		between = "not between"
	}
	column := grammarSQL.Wrap(having.Column)
	if !dbal.IsExpression(having.Values[0]) {
		*bindingOffset = *bindingOffset + having.Offset
	}
	min := grammarSQL.Parameter(having.Values[0], *bindingOffset)

	if !dbal.IsExpression(having.Values[1]) {
		*bindingOffset = *bindingOffset + having.Offset
	}
	max := grammarSQL.Parameter(having.Values[1], *bindingOffset)

	// and `field` between 3 and 5
	// or `field` not between 3 and 5
	return fmt.Sprintf("%s %s %s %s and %s", having.Boolean, column, between, min, max)
}

// CompileOrders Compile the "order by" portions of the query.
func (grammarSQL SQL) CompileOrders(query *dbal.Query, orders []dbal.Order, bindingOffset *int) string {
	if len(orders) == 0 {
		return ""
	}

	clauses := []string{}
	for _, order := range orders {
		if order.SQL != "" {
			clauses = append(clauses, order.SQL)
		} else {
			clauses = append(clauses, fmt.Sprintf("%s %s", grammarSQL.Wrap(order.Column), order.Direction))
		}
	}
	return fmt.Sprintf("order by %s", strings.Join(clauses, ", "))
}

// CompileLimit Compile the "limit" portions of the query.
func (grammarSQL SQL) CompileLimit(query *dbal.Query, limit int, bindingOffset *int) string {
	if limit < 0 {
		return ""
	}
	return fmt.Sprintf("limit %d", limit)
}

// CompileOffset Compile the "offset" portions of the query.
func (grammarSQL SQL) CompileOffset(query *dbal.Query, offset int) string {
	if offset < 0 {
		return ""
	}
	return fmt.Sprintf("offset %d", offset)
}

// CompileLock the lock into SQL.
func (grammarSQL SQL) CompileLock(query *dbal.Query, lock interface{}) string {
	sql, ok := lock.(string)
	if ok {
		return sql
	}
	return ""
}

// CompileWheres Compile an update statement into SQL.
func (grammarSQL SQL) CompileWheres(query *dbal.Query, wheres []dbal.Where, bindingOffset *int) string {

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

// WhereBasic Compile a date based where clause.
func (grammarSQL SQL) WhereBasic(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}

	operator := strings.ReplaceAll(where.Operator, "?", "??")

	return fmt.Sprintf("%s %s %s", grammarSQL.Wrap(where.Column), operator, value)
}

// WhereColumn Compile a where clause comparing two columns.
func (grammarSQL SQL) WhereColumn(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return fmt.Sprintf("%s %s %s", grammarSQL.Wrap(where.First), where.Operator, grammarSQL.Wrap(where.Second))
}

// WhereDate Compile a "where date" clause.
func (grammarSQL SQL) WhereDate(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("date", query, where, bindingOffset)
}

// WhereTime Compile a "where time" clause.
func (grammarSQL SQL) WhereTime(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("time", query, where, bindingOffset)
}

// WhereDay Compile a "where day" clause.
func (grammarSQL SQL) WhereDay(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("day", query, where, bindingOffset)
}

// WhereMonth Compile a "where month" clause.
func (grammarSQL SQL) WhereMonth(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("month", query, where, bindingOffset)
}

// WhereYear Compile a "where year" clause.
func (grammarSQL SQL) WhereYear(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return grammarSQL.WhereDateBased("year", query, where, bindingOffset)
}

// WhereDateBased  Compile a date based where clause.
func (grammarSQL SQL) WhereDateBased(typ string, query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	value := ""
	if !dbal.IsExpression(where.Value) {
		*bindingOffset = *bindingOffset + where.Offset
		value = grammarSQL.Parameter(where.Value, *bindingOffset)
	} else {
		value = where.Value.(dbal.Expression).GetValue()
	}

	return fmt.Sprintf("%s(%s)%s%s", typ, grammarSQL.Wrap(where.Column), where.Operator, value)
}

// WhereNested Compile a nested where clause.
func (grammarSQL SQL) WhereNested(query *dbal.Query, where dbal.Where, bindingOffset *int) string {

	offset := 6 // - where
	if query.IsJoinClause {
		offset = 3 // - on
	}

	sql := grammarSQL.CompileWheres(where.Query, where.Query.Wheres, bindingOffset)
	end := len(sql)
	if end > offset {
		sql = sql[offset:end]
	}
	return fmt.Sprintf("(%s)", sql)
}

// WhereSub Compile a where condition with a sub-select.
func (grammarSQL SQL) WhereSub(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	selectSQL := grammarSQL.CompileSelectOffset(where.Query, bindingOffset)
	return fmt.Sprintf("%s %s (%s)", grammarSQL.Wrap(where.Column), where.Operator, selectSQL)
}

// WhereExists Compile a where (not) exists clause.
func (grammarSQL SQL) WhereExists(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	exists := "exists"
	if where.Not {
		exists = "not exists"
	}
	selectSQL := grammarSQL.CompileSelectOffset(where.Query, bindingOffset)
	return fmt.Sprintf("%s (%s)", exists, selectSQL)
}

// WhereNull Compile a "where null" clause.
func (grammarSQL SQL) WhereNull(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return fmt.Sprintf("%s is null", grammarSQL.Wrap(where.Column))
}

// WhereRaw Compile a raw where clause.
func (grammarSQL SQL) WhereRaw(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return where.SQL
}

// WhereNotnull Compile a "where not null" clause.
func (grammarSQL SQL) WhereNotnull(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	return fmt.Sprintf("%s is not null", grammarSQL.Wrap(where.Column))
}

// WhereBetween Compile a "between" where clause.
func (grammarSQL SQL) WhereBetween(query *dbal.Query, where dbal.Where, bindingOffset *int) string {
	if len(where.Values) != 2 {
		panic(fmt.Errorf("The given values must have two items"))
	}
	between := "between"
	if where.Not {
		between = "not between"
	}
	column := grammarSQL.Wrap(where.Column)
	*bindingOffset = *bindingOffset + where.Offset
	min := grammarSQL.Parameter(where.Values[0], *bindingOffset)
	*bindingOffset = *bindingOffset + 1
	max := grammarSQL.Parameter(where.Values[1], *bindingOffset)
	// `field` between 3 and 5
	// `field` not between 3 and 5
	return fmt.Sprintf("%s %s %s and %s", column, between, min, max)
}

// WhereIn Compile a "where in" clause.
func (grammarSQL SQL) WhereIn(query *dbal.Query, where dbal.Where, bindingOffset *int) string {

	in := "in"
	sql := "false = true"
	if where.Not {
		in = "not in"
		sql = "true = true"
	}

	if where.ValuesIn != nil {
		reflectValues := reflect.ValueOf(where.ValuesIn)
		reflectValues = reflect.Indirect(reflectValues)
		if reflectValues.Kind() == reflect.Slice || reflectValues.Kind() == reflect.Array {
			*bindingOffset = *bindingOffset + where.Offset
			values := []interface{}{}
			for i := 0; i < reflectValues.Len(); i++ {
				values = append(values, reflectValues.Index(i).Interface())
			}

			sql = fmt.Sprintf("%s %s (%s)", grammarSQL.Wrap(where.Column), in, grammarSQL.Parameterize(values, *bindingOffset))
			*bindingOffset = *bindingOffset + len(values)
		} else if _, ok := where.ValuesIn.(dbal.Expression); ok {
			*bindingOffset = *bindingOffset + where.Offset
			sql = fmt.Sprintf("%s %s (%s)", grammarSQL.Wrap(where.Column), in, where.ValuesIn.(dbal.Expression).GetValue())
		}
	}

	return sql
}

// Utils for compiling

// RemoveLeadingBoolean Remove the leading boolean from a statement.
func (grammarSQL SQL) RemoveLeadingBoolean(value string) string {
	value = strings.TrimPrefix(value, "and ")
	value = strings.TrimPrefix(value, "or ")
	return value
}

// Raw make a new expression
func (grammarSQL SQL) Raw(value string) dbal.Expression {
	return dbal.NewExpression(value)
}
