package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// prepareArgs Prepare the value, operator, boolean and offset for a where clause.
func (builder *Builder) prepareArgs(args ...interface{}) (string, interface{}, string, int) {

	var operator string = "="
	var value interface{} = nil
	var boolean string = "and"
	var offset = 1

	// Where("score", 5)
	if len(args) == 1 {
		value = args[0]
		return operator, value, boolean, offset
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

	// Where("vote", ">", 5, "and", 5)
	if len(args) == 4 && reflect.TypeOf(args[3]).Kind() == reflect.Int {
		offset = args[3].(int)
	}

	return operator, value, boolean, offset
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

func (builder *Builder) isOperator(v interface{}) bool {
	switch v.(type) {
	case string:
		return utils.StringHave([]string{"and", "or"}, strings.ToLower(v.(string)))
	default:
		return false
	}
}

// isExpression Determine if the given value is a raw expression.
func (builder *Builder) isExpression(value interface{}) bool {
	switch value.(type) {
	case dbal.Expression:
		return true
	default:
		return false
	}
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

	return builder.isClosure(value) ||
		(kind == reflect.Interface && typ.Name() == "Query") ||
		(kind == reflect.Struct && typ.Name() == "Builder")
}

// Determine if the given operator and value combination is legal.
func (builder *Builder) invalidOperator(operator string) bool {
	return !utils.StringHave(builder.Query.Operators, operator) &&
		!utils.StringHave(builder.Grammar.GetOperators(), operator)
}

// Determine if the given operator is supported.
func (builder *Builder) invalidOperatorAndValue(operator string, value interface{}) bool {
	return value == nil &&
		utils.StringHave(builder.Query.Operators, operator) &&
		utils.StringHave([]string{"=", "<>", "!="}, operator)
}

// Remove all of the expressions from a list of bindings.
func (builder *Builder) cleanBindings(bindings interface{}) []interface{} {

	values, ok := bindings.([]interface{})
	if !ok {
		panic(fmt.Errorf("The input bindings must be the interface slice"))
	}

	for index, value := range values {
		if builder.isExpression(value) {
			values = append(values[:index], values[index+1:]...)
		}
	}

	return values
}

// Get a scalar type value from an unknown type of input.
func (builder *Builder) flattenValue(value interface{}) interface{} {
	values := utils.Flatten(value)
	return values[0]
}
