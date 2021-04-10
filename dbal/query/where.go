package query

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Where Add a basic where clause to the query.
func (builder *Builder) Where(column interface{}, args ...interface{}) Query {

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, value, boolean, offset := builder.prepareWhereArgs(args...)
	whereType := "basic"

	// If the columns is actually a Closure instance, we will assume the developer
	// wants to begin a nested where statement which is wrapped in parenthesis.
	// We'll add that Closure to the query then return back out immediately.
	if builder.isClosure(column) && (len(args) == 0 || (len(args) >= 1 && builder.isBoolean(args[0]))) {
		boolean = "and"
		if len(args) == 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
			boolean = args[0].(string)
		}
		whereType = "nested"

		// If the column is a Closure instance and there is an operator value, we will
		// assume the developer wants to run a subquery and then compare the result
		// of that subquery with the given value that was provided to the method.
	} else if builder.isQueryable(column) && len(args) >= 1 {
		whereType = "sub"
	}

	return builder.where(column, operator, value, boolean, offset, whereType)
}

// OrWhere Add an "or where" clause to the query.
func (builder *Builder) OrWhere(column interface{}, args ...interface{}) Query {

	if builder.isClosure(column) && (len(args) == 0 || (len(args) >= 1 && builder.isBoolean(args[0]))) {
		return builder.Where(column, "or")
	}

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, value, _, offset := builder.prepareWhereArgs(args...)
	return builder.Where(column, operator, value, "or", offset)

}

// Where Add a basic where clause to the query.
func (builder *Builder) where(column interface{}, operator string, value interface{}, boolean string, offset int, whereType string) Query {

	columnKind := reflect.TypeOf(column).Kind()

	// Where([][]interface{}{ {"score", ">", 64.56},{"vote", 10}})
	// If the column is an array, we will assume it is an array of key-value pairs
	// and can add them each as a where clause. We will maintain the boolean we
	// received when the method was called and pass it into the Wheres attribute.
	if columnKind == reflect.Array || columnKind == reflect.Slice || columnKind == reflect.Map {
		builder.addArrayOfWheres(column, boolean)
		return builder
	}

	// Where( func(qb Query){  qb.Where("name", "Ken")... })
	// Where( func(qb Query){  qb.Where("name", "Ken")... }, "and")
	// Where( func(qb Query){  qb.Where("name", "Ken")... }, "or")
	if whereType == "nested" {
		callback := column.(func(Query))
		builder.whereNested(callback, boolean)
		return builder
	}

	// Where( func(qb Query){ qb.Where("name", "Ken")... }, 5)
	// Where( func(qb Query){ qb.Where("name", "Ken")... }, 5, "or")
	// Where( func(qb Query){ qb.Where("name", "Ken")... }, ">", 5)
	// Where( func(qb Query){ qb.Where("name", "Ken")... }, ">", 5, "or")
	if whereType == "sub" {
		sub, bindings, whereOffset := builder.createSub(column)
		segment := builder.parseSub(sub)
		builder.Query.AddBinding("where", bindings)
		builder.Where(dbal.Raw(fmt.Sprintf("(%s)", segment)), operator, value, boolean, whereOffset)
		return builder
	}

	// Where("vote", '>', func(sub Query) {
	// 		sub.From("table_test_where").
	// 			Where("score", ">", 5).
	// 	 		Sum()
	// })
	// If the value is a Closure, it means the developer is performing an entire
	// sub-select within the query and we will need to compile the sub-select
	// within the where clause to get the appropriate query record results.
	if columnKind == reflect.String && builder.isClosure(value) {
		callback := value.(func(Query))
		return builder.whereSub(column.(string), operator, callback, boolean)
	}

	// If the value is "null", we will just assume the developer wants to add a
	// where null clause to the query. So, we will allow a short-cut here to
	// that method for convenience so the developer doesn't have to check.
	if utils.IsNil(value) {
		return builder.WhereNull(column, boolean, operator != "=")
	}

	queryType := "basic"

	// If the column is making a JSON reference we'll check to see if the value
	// is a boolean. If it is, we'll add the raw boolean string as an actual
	// value to the query to ensure this is properly handled by the query.
	//  if (Str::contains($column, '->') && is_bool($value)) {

	// Where("email", "like", "%@yao.run")
	// Now that we are working with just a simple query we can put the elements
	// in our array and add the query binding to our array of bindings that
	// will be bound to each SQL statements when it is finally executed.
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:     queryType,
		Column:   column,
		Operator: operator,
		Boolean:  boolean,
		Value:    value,
		Offset:   offset,
	})
	if !builder.isExpression(value) {
		builder.Query.AddBinding("where", builder.flattenValue(value))
	}
	return builder
}

// WhereColumn Add a "where" clause comparing two columns to the query.
func (builder *Builder) WhereColumn(first interface{}, args ...interface{}) Query {

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, second, _, offset := builder.prepareWhereArgs(args...)

	return builder.whereColumn(first, operator, second, "and", offset)
}

// OrWhereColumn Add an "or where" clause comparing two columns to the query.
func (builder *Builder) OrWhereColumn(first interface{}, args ...interface{}) Query {

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, second, _, offset := builder.prepareWhereArgs(args...)

	return builder.whereColumn(first, operator, second, "or", offset)
}

// WhereColumn Add a "where" clause comparing two columns to the query.
func (builder *Builder) whereColumn(first interface{}, operator string, second interface{}, boolean string, offset int) Query {

	columnKind := reflect.TypeOf(first).Kind()

	// Where([][]interface{}{ {"score", ">", 64.56},{"vote", 10}})
	// If the column is an array, we will assume it is an array of key-value pairs
	// and can add them each as a where clause. We will maintain the boolean we
	// received when the method was called and pass it into the Wheres attribute.
	if columnKind == reflect.Array || columnKind == reflect.Slice {
		builder.addArrayOfWheres(first, boolean, true)
		return builder
	}

	// Finally, we will add this where clause into this array of clauses that we
	// are building for the query. All of them will be compiled via a grammar
	// once the query is about to be executed and run against the database.
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:     "column",
		Second:   second,
		First:    first,
		Operator: operator,
		Boolean:  boolean,
		Offset:   offset,
	})

	return builder
}

// Where([][]interface{}{ {"score", ">", 64.56},{"vote", 10},})
// addArrayOfWheres Add an array of where clauses to the query.
func (builder *Builder) addArrayOfWheres(inputColumns interface{}, boolean string, whereColumn ...bool) *Builder {

	switch inputColumns.(type) {
	case map[string]interface{}, xun.R:
		columns, ok := inputColumns.(map[string]interface{})
		if !ok {
			columns = inputColumns.(xun.R).ToMap()
		}
		return builder.whereNested(func(qb Query) {
			for column, value := range columns {
				qb.Where(column, "=", value, "and", 1)
			}
		}, boolean)

	case [][]interface{}:
		columns := inputColumns.([][]interface{})
		return builder.whereNested(func(qb Query) {
			for _, args := range columns {
				if len(args) > 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
					column := args[0].(string)
					operator, value, boolean, offset := builder.prepareWhereArgs(args[1:]...)
					if len(whereColumn) > 0 && whereColumn[0] {
						qb.WhereColumn(column, operator, value, boolean, offset)
					} else {
						qb.Where(column, operator, value, boolean, offset)
					}
				}
			}
		}, boolean)
	}

	return builder
}

// createSub Creates a subquery and parse it.
func (builder *Builder) createSub(subquery interface{}) (interface{}, []interface{}, int) {
	if builder.isClosure(subquery) {
		callback := subquery.(func(qb Query))
		qb := builder.forSubQuery()
		callback(qb)
		subquery = qb
	}
	return builder.makeSub(subquery)
}

// make the subquery and bindings.
func (builder *Builder) makeSub(subquery interface{}) (interface{}, []interface{}, int) {
	switch subquery.(type) {
	case *Builder:
		qb := builder.prependDatabaseNameIfCrossDatabaseQuery(subquery.(*Builder))
		offset := len(builder.GetBindings())
		bindings := qb.GetBindings()
		whereOffset := offset + len(utils.Flatten(bindings))
		qb.Query.BindingOffset = offset
		return qb.Query, bindings, whereOffset
	case dbal.Expression:
		return subquery.(dbal.Expression).GetValue(), []interface{}{}, 1
	case string:
		return subquery.(string), []interface{}{}, 1
	}

	panic(fmt.Errorf("a subquery must be a query builder instance, a Closure, or a string"))
}

// prependDatabaseNameIfCrossDatabaseQuery Prepend the database name if the given query is on another database.
func (builder *Builder) prependDatabaseNameIfCrossDatabaseQuery(subquery *Builder) *Builder {
	// if (strpos($query->from, $databaseName) !== 0 && strpos($query->from, '.') === false) {
	// 	$query->from($databaseName.'.'.$query->from);
	// }
	return subquery
}

// whereSub Add a full sub-select to the query.
func (builder *Builder) whereSub(column string, operator string, callback func(qb Query), boolean string) *Builder {
	new := builder.forSubQuery()
	callback(new)
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:     "sub",
		Column:   column,
		Operator: operator,
		Query:    new.Query,
		Boolean:  boolean,
	})
	builder.Query.AddBinding("where", new.Query.Bindings["where"])
	return builder
}

// forSubQuery Create a new query instance for a sub-query.
func (builder *Builder) forSubQuery() *Builder {
	new := builder.new()
	return new
}

// whereNested  Add a nested where statement to the query.
func (builder *Builder) whereNested(callback func(qb Query), boolean string) *Builder {
	new := builder.forNestedWhere()
	callback(new)
	return builder.addNestedWhereQuery(new.Query, boolean)
}

// Add another query builder as a nested where to the query builder.
func (builder *Builder) addNestedWhereQuery(query *dbal.Query, boolean string) *Builder {

	if len(query.Wheres) > 0 {
		builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
			Type:    "nested",
			Query:   query,
			Boolean: boolean,
		})
		builder.Query.AddBinding("where", query.Bindings["where"])
	}

	return builder
}

// forNestedWhere Create a new query instance for nested where condition.
func (builder *Builder) forNestedWhere() *Builder {
	new := builder.new()
	new.Query.From = builder.Query.From
	new.Query.IsJoinClause = builder.Query.IsJoinClause
	return new
}

// WhereNull Add a "where null" clause to the query.
func (builder *Builder) WhereNull(column interface{}, args ...interface{}) Query {

	not := false
	boolean := "and"
	typ := "null"
	if len(args) > 0 && builder.isBoolean(args[0]) {
		boolean = args[0].(string)
	}
	if len(args) > 1 && reflect.TypeOf(args[1]).Kind() == reflect.Bool {
		not = args[1].(bool)
	}
	if not {
		typ = "notnull"
	}

	reflectColumn := reflect.ValueOf(column)
	columnKind := reflectColumn.Kind()

	columns := []interface{}{}
	if columnKind == reflect.Array || columnKind == reflect.Slice {
		reflectColumn = reflect.Indirect(reflectColumn)
		for i := 0; i < reflectColumn.Len(); i++ {
			if reflectColumn.Index(i).Kind() == reflect.String {
				columns = append(columns, reflectColumn.Index(i).String())
			}
		}
	} else if columnKind == reflect.String {
		columns = append(columns, reflectColumn.String())
	}

	for _, col := range columns {
		builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
			Column:  col,
			Type:    typ,
			Boolean: boolean,
		})
	}
	return builder
}

// OrWhereNull Add an "or where null" clause to the query.
func (builder *Builder) OrWhereNull(column interface{}) Query {
	return builder.WhereNull(column, "or", false)
}

// WhereNotNull Add a "where not null" clause to the query.
func (builder *Builder) WhereNotNull(column interface{}, args ...interface{}) Query {
	boolean, _, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereNull(column, boolean, true)
}

// OrWhereNotNull Add an "or where not null" clause to the query.
func (builder *Builder) OrWhereNotNull(column interface{}) Query {
	return builder.WhereNull(column, "or", true)
}

// WhereRaw Add a basic where clause to the query.
func (builder *Builder) WhereRaw(sql string, bindings ...interface{}) Query {
	return builder.whereRaw(sql, bindings, "and")
}

// OrWhereRaw Add an "or where" clause to the query.
func (builder *Builder) OrWhereRaw(sql string, bindings ...interface{}) Query {
	return builder.whereRaw(sql, bindings, "or")
}

func (builder *Builder) whereRaw(sql string, bindings []interface{}, boolean string) Query {
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:    "raw",
		SQL:     sql,
		Boolean: boolean,
	})
	builder.Query.AddBinding("where", bindings)
	return builder
}

// WhereBetween Add a where between statement to the query.
func (builder *Builder) WhereBetween(column interface{}, values interface{}) Query {
	return builder.whereBetween(column, values, "and", false)
}

// OrWhereBetween Add an or where between statement to the query.
func (builder *Builder) OrWhereBetween(column interface{}, values interface{}) Query {
	return builder.whereBetween(column, values, "or", false)
}

// WhereNotBetween Add a where not between statement to the query.
func (builder *Builder) WhereNotBetween(column interface{}, values interface{}) Query {
	return builder.whereBetween(column, values, "and", true)
}

// OrWhereNotBetween Add an or where not between statement using columns to the query.
func (builder *Builder) OrWhereNotBetween(column interface{}, values interface{}) Query {
	return builder.whereBetween(column, values, "or", true)
}

// whereBetween Add a where between statement to the query.
func (builder *Builder) whereBetween(column interface{}, values interface{}, boolean string, not bool) Query {

	betweenOffset := 1

	betweenValues := builder.cleanBindings(values)

	if len(betweenValues) > 2 {
		betweenValues = betweenValues[0:2]
	}

	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:    "between",
		Column:  column,
		Values:  betweenValues,
		Boolean: boolean,
		Not:     not,
		Offset:  betweenOffset,
	})

	builder.Query.AddBinding("where", betweenValues)
	return builder
}

// WhereIn Add a "where in" clause to the query.
func (builder *Builder) WhereIn(column interface{}, values interface{}) Query {
	return builder.whereIn(column, values, "and", false)
}

// OrWhereIn Add an "or where in" clause to the query.
func (builder *Builder) OrWhereIn(column interface{}, values interface{}) Query {
	return builder.whereIn(column, values, "or", false)
}

// WhereNotIn Add a "where not in" clause to the query.
func (builder *Builder) WhereNotIn(column interface{}, values interface{}) Query {
	return builder.whereIn(column, values, "and", true)
}

// OrWhereNotIn Add an "or where not in" clause to the query.
func (builder *Builder) OrWhereNotIn(column interface{}, values interface{}) Query {
	return builder.whereIn(column, values, "or", true)
}

// whereIn Add a "where in" clause to the query.
func (builder *Builder) whereIn(column interface{}, values interface{}, boolean string, not bool) Query {

	inOffset := 0

	// If the value is a query builder instance we will assume the developer wants to
	// look for any values that exists within this given query. So we will add the
	// query accordingly so that this query is properly executed when it is run.
	if builder.isQueryable(values) {
		sub, bindings, offset := builder.createSub(values)
		segment := builder.parseSub(sub)
		values = dbal.Raw(segment)
		builder.Query.AddBinding("where", bindings)
		inOffset = offset
	}

	where := dbal.Where{
		Type:     "in",
		ValuesIn: values,
		Column:   column,
		Offset:   inOffset,
		Not:      not,
		Boolean:  boolean,
	}

	builder.Query.Wheres = append(builder.Query.Wheres, where)

	// Finally we'll add a binding for each values unless that value is an expression
	// in which case we will just skip over it since it will be the query as a raw
	// string and not as a parameterized place-holder to be replaced by the PDO.
	builder.Query.AddBinding("where", builder.cleanBindings(values))

	return builder
}

// WhereExists Add an exists clause to the query.
func (builder *Builder) WhereExists(closure func(qb Query)) Query {
	return builder.whereExists(closure, "and", false)
}

// OrWhereExists Add an or exists clause to the query.
func (builder *Builder) OrWhereExists(closure func(qb Query)) Query {
	return builder.whereExists(closure, "or", false)
}

// WhereNotExists  Add a where not exists clause to the query.
func (builder *Builder) WhereNotExists(closure func(qb Query)) Query {
	return builder.whereExists(closure, "and", true)
}

// OrWhereNotExists Add a where not exists clause to the query.
func (builder *Builder) OrWhereNotExists(closure func(qb Query)) Query {
	return builder.whereExists(closure, "or", true)
}

// WhereExists Add an exists clause to the query.
func (builder *Builder) whereExists(closure func(qb Query), boolean string, not bool) Query {
	new := builder.forSubQuery()
	closure(new)
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:    "exists",
		Not:     not,
		Boolean: boolean,
		Query:   new.Query,
	})
	builder.Query.AddBinding("where", new.GetBindings())
	return builder
}

// WhereDate Add a "where date" statement to the query.
func (builder *Builder) WhereDate(column interface{}, args ...interface{}) Query {
	operator, value, boolean, _ := builder.prepareWhereArgs(args...)
	// date Format
	if !builder.isExpression(value) {
		value = xun.MakeTime(value).MustToTime().Format("2006-01-02")
	}
	return builder.whereDateTime("date", column, operator, value, boolean)
}

// OrWhereDate Add an "or where date" statement to the query.
func (builder *Builder) OrWhereDate(column interface{}, args ...interface{}) Query {
	operator, value, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereDate(column, operator, value, "or")
}

// WhereTime Add a "where time" statement to the query.
func (builder *Builder) WhereTime(column interface{}, args ...interface{}) Query {
	operator, value, boolean, _ := builder.prepareWhereArgs(args...)
	// date Format
	if !builder.isExpression(value) {
		value = xun.MakeTime(value).MustToTime().Format("15:04:05")
	}
	return builder.whereDateTime("time", column, operator, value, boolean)
}

// OrWhereTime Add an "or where time" statement to the query.
func (builder *Builder) OrWhereTime(column interface{}, args ...interface{}) Query {
	operator, value, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereTime(column, operator, value, "or")
}

// WhereYear Add a "where year" statement to the query.
func (builder *Builder) WhereYear(column interface{}, args ...interface{}) Query {
	operator, value, boolean, _ := builder.prepareWhereArgs(args...)
	// date Format
	if !builder.isExpression(value) {
		value = xun.MakeTime(value).MustToTime().Year()
	}
	return builder.whereDateTime("year", column, operator, value, boolean)
}

// OrWhereYear Add an "or where year" statement to the query.
func (builder *Builder) OrWhereYear(column interface{}, args ...interface{}) Query {
	operator, value, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereYear(column, operator, value, "or")
}

// WhereMonth Add a "where month" statement to the query.
func (builder *Builder) WhereMonth(column interface{}, args ...interface{}) Query {
	operator, value, boolean, _ := builder.prepareWhereArgs(args...)
	// date Format
	if !builder.isExpression(value) {
		value = xun.MakeTime(value).MustToTime().Month()
	}
	return builder.whereDateTime("month", column, operator, value, boolean)
}

// OrWhereMonth Add an "or where month" statement to the query.
func (builder *Builder) OrWhereMonth(column interface{}, args ...interface{}) Query {
	operator, value, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereMonth(column, operator, value, "or")
}

// WhereDay Add a "where day" statement to the query.
func (builder *Builder) WhereDay(column interface{}, args ...interface{}) Query {
	operator, value, boolean, _ := builder.prepareWhereArgs(args...)
	// date Format
	if !builder.isExpression(value) {
		value = xun.MakeTime(value).MustToTime().Day()
	}
	return builder.whereDateTime("day", column, operator, value, boolean)
}

// OrWhereDay Add an "or where day" statement to the query.
func (builder *Builder) OrWhereDay(column interface{}, args ...interface{}) Query {
	operator, value, _, _ := builder.prepareWhereArgs(args...)
	return builder.WhereDay(column, operator, value, "or")
}

// whereDateTime Add a date based (date, time, year, month, day ) statement to the query.
func (builder *Builder) whereDateTime(dateType string, column interface{}, operator string, value interface{}, boolean string) Query {

	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:     dateType,
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
		Offset:   1,
	})

	if !builder.isExpression(value) {
		builder.Query.AddBinding("where", value)
	}

	return builder
}

// WhereJSONContains Add a "where JSON contains" clause to the query.
func (builder *Builder) WhereJSONContains() {
}

// OrWhereJSONContains Add an "or where JSON contains" clause to the query.
func (builder *Builder) OrWhereJSONContains() {
}

// WhereJSONDoesntContain Add a "where JSON not contains" clause to the query.
func (builder *Builder) WhereJSONDoesntContain() {
}

// OrWhereJSONDoesntContain Add an "or where JSON not contains" clause to the query.
func (builder *Builder) OrWhereJSONDoesntContain() {
}

// WhereJSONLength Add a "where JSON length" clause to the query.
func (builder *Builder) WhereJSONLength() {
}

// OrWhereJSONLength Add an "or where JSON length" clause to the query.
func (builder *Builder) OrWhereJSONLength() {
}
