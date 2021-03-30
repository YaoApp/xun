package query

import (
	"reflect"

	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Where Add a basic where clause to the query.
func (builder *Builder) Where(column interface{}, args ...interface{}) Query {

	typ := reflect.TypeOf(column)
	columnKind := typ.Kind()
	queryType := "basic"

	// Where([][]interface{}{ {"score", ">", 64.56},{"vote", 10}})
	// If the column is an array, we will assume it is an array of key-value pairs
	// and can add them each as a where clause. We will maintain the boolean we
	// received when the method was called and pass it into the Wheres attribute.
	if columnKind == reflect.Array || columnKind == reflect.Slice {
		builder.addArrayOfWheres(queryType, column)
		return builder
	}

	// Here we will make some assumptions about the operator. If only 2 values are
	// passed to the method, we will assume that the operator is an equals sign
	// and keep going. Otherwise, we'll require the operator to be passed in.
	operator, value, boolean := builder.wherePrepare(args...)

	// Where( func(qb Query){  qb.Where("name", "Ken")... })
	// Where( func(qb Query){  qb.Where("name", "Ken")... }, "and")
	// Where( func(qb Query){  qb.Where("name", "Ken")... }, "or")
	// If the columns is actually a Closure instance, we will assume the developer
	// wants to begin a nested where statement which is wrapped in parenthesis.
	// We'll add that Closure to the query then return back out immediately.
	if builder.isClosure(column) && len(args) == 0 {
		callback := column.(func(Query))
		boolean := "and"
		if len(args) == 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
			boolean = args[0].(string)
		}
		builder.whereNested(callback, boolean)
		return builder
	}

	// If the column is a Closure instance and there is an operator value, we will
	// assume the developer wants to run a subquery and then compare the result
	// of that subquery with the given value that was provided to the method.
	// if ($this->isQueryable($column) && ! is_null($operator)) {

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

	// // If the column is an string, we willpass it into the Wheres attribute.
	// if columnKind == reflect.String {
	// 	builder.addWhere(queryType, column.(string), operator, value, boolean)
	// 	return builder
	// }

	queryType = "basic"

	// If the column is making a JSON reference we'll check to see if the value
	// is a boolean. If it is, we'll add the raw boolean string as an actual
	// value to the query to ensure this is properly handled by the query.
	//  if (Str::contains($column, '->') && is_bool($value)) {

	// Where("email", "like", "%@yao.run")
	// Now that we are working with just a simple query we can put the elements
	// in our array and add the query binding to our array of bindings that
	// will be bound to each SQL statements when it is finally executed.
	if columnKind == reflect.String {
		builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
			Type:     queryType,
			Column:   column.(string),
			Operator: operator,
			Boolean:  boolean,
			Value:    value,
		})
		builder.Query.AddBinding("where", builder.flattenValue(value))
	}
	return builder
}

// Where([][]interface{}{ {"score", ">", 64.56},{"vote", 10},})
// addArrayOfWheres Add an array of where clauses to the query.
func (builder *Builder) addArrayOfWheres(typ string, inputColumns interface{}) {

	switch inputColumns.(type) {
	case [][]interface{}:
		reflectValue := reflect.ValueOf(inputColumns)
		reflectValue = reflect.Indirect(reflectValue)
		for i := 0; i < reflectValue.Len(); i++ {
			args := reflectValue.Index(i).Interface().([]interface{})
			if len(args) > 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
				column := args[0].(string)
				operator, value, boolean := builder.wherePrepare(args[1:]...)
				builder.addWhere(typ, column, operator, value, boolean)
			}
		}
		break
	}

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
	return new
}

func (builder *Builder) addWhere(typ string, column string, operator string, value interface{}, boolean string) {
	builder.Query.Wheres = append(builder.Query.Wheres, dbal.Where{
		Type:     typ,
		Column:   column,
		Operator: operator,
		Boolean:  boolean,
		Wheres:   []dbal.Where{},
	})
}

// wherePrepare Prepare the value, operator and boolean for a where clause.
func (builder *Builder) wherePrepare(args ...interface{}) (string, interface{}, string) {
	var operator string = "="
	var value interface{} = nil
	var boolean string = "and"

	// Where("score", 5)
	if len(args) == 1 {
		value = args[0]
		return operator, value, boolean
	}

	// Where("vote", ">", 5)
	if len(args) >= 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
		operator = args[0].(string)
	}

	if len(args) >= 2 {
		value = args[1]
	}

	// Where("vote", ">", 5, "and")
	if len(args) >= 3 && reflect.TypeOf(args[2]).Kind() == reflect.String {
		boolean = args[2].(string)
	}

	return operator, value, boolean
}

func (builder *Builder) isClosure(v interface{}) bool {
	if v == nil {
		return false
	}
	typ := reflect.TypeOf(v)
	return typ.Kind() == reflect.Func &&
		typ.NumOut() == 0 &&
		typ.NumIn() == 1 &&
		typ.In(0).Kind() == reflect.Interface
}

// Determine if the value is a query builder instance or a Closure.
func (builder *Builder) isQueryable(value interface{}) bool {
	typ := reflect.TypeOf(value)
	kind := typ.Kind()
	if kind == reflect.Ptr {
		reflectValue := reflect.Indirect(reflect.ValueOf(value))
		typ = reflectValue.Type()
		kind = typ.Kind()
	}
	return (kind == reflect.Interface && typ.Name() == "Query") ||
		(kind == reflect.Struct && typ.Name() == "Builder")
}

func (builder *Builder) invalidOperator(operator string) bool {
	return !utils.StringHave(builder.Query.Operators, operator) &&
		!utils.StringHave(builder.Grammar.GetOperators(), operator)
}

func (builder *Builder) invalidOperatorAndValue(operator string, value interface{}) bool {
	return value == nil &&
		utils.StringHave(builder.Query.Operators, operator) &&
		utils.StringHave([]string{"=", "<>", "!="}, operator)
}

func (builder *Builder) flattenValue(value interface{}) interface{} {
	values := utils.Flatten(value)
	return values[0]
}

// OrWhere Add an "or where" clause to the query.
func (builder *Builder) OrWhere() {
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

// WhereBetween Add a where between statement to the query.
func (builder *Builder) WhereBetween() {
}

// OrWhereBetween Add an or where between statement to the query.
func (builder *Builder) OrWhereBetween() {
}

// WhereNotBetween Add a where not between statement to the query.
func (builder *Builder) WhereNotBetween() {
}

// OrWhereNotBetween Add an or where not between statement using columns to the query.
func (builder *Builder) OrWhereNotBetween() {
}

// WhereIn Add a "where in" clause to the query.
func (builder *Builder) WhereIn() {
}

// OrWhereIn Add an "or where in" clause to the query.
func (builder *Builder) OrWhereIn() {
}

// WhereNotIn Add a "where not in" clause to the query.
func (builder *Builder) WhereNotIn() {
}

// OrWhereNotIn Add an "or where not in" clause to the query.
func (builder *Builder) OrWhereNotIn() {
}

// WhereNull Add a "where null" clause to the query.
func (builder *Builder) WhereNull(column interface{}, args ...interface{}) *Builder {
	boolean, not, _ := builder.wherePrepare(args...)
	typ := "null"
	if !utils.IsNil(not) && reflect.TypeOf(not).Kind() == reflect.Bool {
		if reflect.ValueOf(not).Bool() {
			typ = "notnull"
		}
	}

	reflectColumn := reflect.ValueOf(column)
	columnKind := reflectColumn.Kind()

	columns := []string{}
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
func (builder *Builder) OrWhereNull(column interface{}) *Builder {
	return builder.WhereNull(column, "or")
}

// WhereNull Add a "where not null" clause to the query.
func (builder *Builder) whereNotNull(column interface{}, args ...interface{}) *Builder {
	boolean, _, _ := builder.wherePrepare(args...)
	return builder.WhereNull(column, boolean, true)
}

// OrWhereNotNull Add an "or where not null" clause to the query.
func (builder *Builder) OrWhereNotNull(column interface{}) *Builder {
	return builder.WhereNull(column, "or", true)
}

// WhereDate Add a "where date" statement to the query.
func (builder *Builder) WhereDate() {
}

// OrWhereDate Add an "or where date" statement to the query.
func (builder *Builder) OrWhereDate() {
}

// WhereYear Add a "where year" statement to the query.
func (builder *Builder) WhereYear() {
}

// OrWhereYear Add an "or where year" statement to the query.
func (builder *Builder) OrWhereYear() {
}

// WhereMonth Add a "where month" statement to the query.
func (builder *Builder) WhereMonth() {
}

// OrWhereMonth Add an "or where month" statement to the query.
func (builder *Builder) OrWhereMonth() {
}

// WhereDay Add a "where day" statement to the query.
func (builder *Builder) WhereDay() {
}

// OrWhereDay Add an "or where day" statement to the query.
func (builder *Builder) OrWhereDay() {
}

// WhereTime Add a "where time" statement to the query.
func (builder *Builder) WhereTime() {
}

// OrWhereTime Add an "or where time" statement to the query.
func (builder *Builder) OrWhereTime() {
}

// WhereColumn Add a "where" clause comparing two columns to the query.
func (builder *Builder) WhereColumn() {
}

// OrWhereColumn Add an "or where" clause comparing two columns to the query.
func (builder *Builder) OrWhereColumn() {
}

// WhereExists Add an exists clause to the query.
func (builder *Builder) WhereExists() {
}

// OrWhereExists Add an or exists clause to the query.
func (builder *Builder) OrWhereExists() {
}

// WhereNotExists  Add a where not exists clause to the query.
func (builder *Builder) WhereNotExists() {
}

// OrWhereNotExists Add a where not exists clause to the query.
func (builder *Builder) OrWhereNotExists() {
}

// WhereRaw Add a basic where clause to the query.
func (builder *Builder) WhereRaw() {
}

// OrWhereRaw Add an "or where" clause to the query.
func (builder *Builder) OrWhereRaw() {
}
