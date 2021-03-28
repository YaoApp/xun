package query

import (
	"reflect"
)

// Where Add a basic where clause to the query.
func (builder *Builder) Where(column interface{}, args ...interface{}) Query {

	typ := reflect.TypeOf(column)
	kind := typ.Kind()
	queryType := "where"

	operator, value, boolean := builder.getDefaultsOfWhere(args...)

	// If the value is a Closure, it means the developer is performing an entire
	// sub-select within the query and we will need to compile the sub-select
	// within the where clause to get the appropriate query record results.
	if builder.isClosure(value) {
		queryType = "sub"
		cb := value.(func(Query))
		valueBuilder := builder.new()
		cb(valueBuilder)
		value = valueBuilder
	}

	// Where("email", "like", "%@yao.run")
	// If the column is an string, we willpass it into the Wheres attribute.
	if kind == reflect.String {
		builder.addWhere(&builder.Attr.Wheres, queryType, column.(string), operator, value, boolean)
		return builder
	}

	// Where([]string{"score", "vote"}, ">", 5)
	// Where([][]interface{}{
	// 	{"score", ">", 64.56},
	// 	{"vote", 10},
	// })
	// If the column is an array, we will assume it is an array of key-value pairs
	// and can add them each as a where clause. We will maintain the boolean we
	// received when the method was called and pass it into the Wheres attribute.
	if kind == reflect.Array || kind == reflect.Slice {
		builder.addWheres(&builder.Attr.Wheres, queryType, column, operator, value, boolean)
		return builder
	}

	// Where( func(qb Query){
	//	  qb.Where("name", "Ken")
	//    ...
	// })
	// If the columns is actually a Closure instance, we will assume the developer
	// wants to begin a nested where statement which is wrapped in parenthesis.
	// We'll add that Closure to the query then return back out immediately.
	if builder.isClosure(column) {

		cb := column.(func(Query))
		newBuilder := builder.new()
		cb(newBuilder)

		boolean := "and"
		if len(args) == 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
			boolean = args[0].(string)
		}

		builder.addNestedWhere(&builder.Attr.Wheres, newBuilder, boolean)
		return builder
	}

	return builder
}

func (builder *Builder) addWheres(wheres *[]Where, typ string, inputColumns interface{}, operator string, value interface{}, boolean string) {

	reflectValue := reflect.ValueOf(inputColumns)
	reflectValue = reflect.Indirect(reflectValue)

	switch inputColumns.(type) {

	// Where([]string{"score", "vote"}, ">", 5)
	case []string:
		for i := 0; i < reflectValue.Len(); i++ {
			column := reflectValue.Index(i)
			if column.Kind() == reflect.String {
				builder.addWhere(wheres, typ, column.String(), operator, value, boolean)
			}
		}
		break

	// Where([][]interface{}{
	// 	{"score", ">", 64.56},
	// 	{"vote", 10},
	// })
	case [][]interface{}:
		for i := 0; i < reflectValue.Len(); i++ {
			args := reflectValue.Index(i).Interface().([]interface{})
			if len(args) > 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
				column := args[0].(string)
				operator, value, boolean := builder.getDefaultsOfWhere(args[1:]...)
				builder.addWhere(wheres, typ, column, operator, value, boolean)
			}
		}
		break
	}

}

func (builder *Builder) addWhere(wheres *[]Where, typ string, column string, operator string, value interface{}, boolean string) {
	*wheres = append(*wheres, Where{
		Type:     typ,
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
		Wheres:   []Where{},
	})
}

func (builder *Builder) addNestedWhere(wheres *[]Where, child *Builder, boolean string) {
	*wheres = append(*wheres, Where{Type: "nested", Query: child, Boolean: boolean})
}

func (builder *Builder) getDefaultsOfWhere(args ...interface{}) (string, interface{}, string) {
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
func (builder *Builder) WhereNull() {
}

// OrWhereNull Add an "or where null" clause to the query.
func (builder *Builder) OrWhereNull() {
}

// WhereNull Add a "where not null" clause to the query.
func (builder *Builder) whereNotNull() {
}

// OrWhereNotNull Add an "or where not null" clause to the query.
func (builder *Builder) OrWhereNotNull() {
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
