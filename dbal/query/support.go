package query

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// prepareWhereArgs Prepare the value, operator, boolean and offset for a where clause.
func (builder *Builder) prepareWhereArgs(args ...interface{}) (string, interface{}, string, int) {

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

// prepareInsertValues prepare the insert values
func (builder *Builder) prepareInsertValues(v interface{}, columns ...interface{}) ([]interface{}, [][]interface{}) {

	if _, ok := v.([][]interface{}); len(columns) > 0 && ok {
		columns = builder.prepareColumns(columns...)
		return columns, v.([][]interface{})
	}

	values := xun.MakeRows(v)
	columns = values[0].Keys()
	insertValues := [][]interface{}{}
	for _, row := range values {
		insertValue := []interface{}{}
		for _, column := range columns {
			insertValue = append(insertValue, row.Get(column))
		}
		insertValues = append(insertValues, insertValue)
	}
	return columns, insertValues
}

// prepareColumns parepare the select columns
// Select("field1", "field2")
// Select("field1", "field2 as f2")
// Select("field1", dbal.Raw("Count(id) as v"))
// Select("field1,field2")
// Select([]string{"field1", "field2"})
func (builder *Builder) prepareColumns(v ...interface{}) []interface{} {

	// columns  "field1,field2", []string{"field1", "field2"}
	if len(v) == 1 {
		col, ok := v[0].(string)
		if ok && strings.Contains(col, ",") {
			cols := strings.Split(col, ",")
			columns := []interface{}{}
			for _, col := range cols {
				columns = append(columns, strings.Trim(col, " "))
			}
			return columns
		} else if !ok {
			reflectValue := reflect.ValueOf(v[0])
			kind := reflectValue.Kind()
			columns := []interface{}{}
			if kind == reflect.Array || kind == reflect.Slice {
				if reflectValue.Len() == 1 {
					col, ok := reflectValue.Index(0).Interface().(string)
					if ok && strings.Contains(col, ",") {
						return builder.prepareColumns(reflectValue.Index(0).Interface())
					}
				}
				for i := 0; i < reflectValue.Len(); i++ {
					columns = append(columns, reflectValue.Index(i).Interface())
				}
				return columns
			}
		}
	}
	return v
}

// Parse the subquery into SQL and bindings.
func (builder *Builder) parseSub(sub interface{}) string {
	switch sub.(type) {
	case *Builder:
		qb := sub.(*Builder)
		offset := qb.Query.BindingOffset
		return qb.Grammar.CompileSelectOffset(qb.Query, &offset)
	case *dbal.Query:
		query := sub.(*dbal.Query)
		offset := query.BindingOffset
		return builder.Grammar.CompileSelectOffset(query, &offset)
	case dbal.Expression:
		return sub.(dbal.Expression).GetValue()
	case string:
		return sub.(string)
	}
	panic(fmt.Errorf("a subquery must be a query builder instance, a Closure, or a string"))
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

func (builder *Builder) isBoolean(v interface{}) bool {
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
	return !utils.StringHave(dbal.Operators, operator) &&
		!utils.StringHave(builder.Grammar.GetOperators(), operator)
}

// Determine if the given operator is supported.
func (builder *Builder) invalidOperatorAndValue(operator string, value interface{}) bool {
	return value == nil &&
		utils.StringHave(dbal.Operators, operator) &&
		utils.StringHave([]string{"=", "<>", "!="}, operator)
}

// Remove all of the expressions from a list of bindings.
func (builder *Builder) cleanBindings(bindings interface{}) []interface{} {
	values := []interface{}{}
	reflectValues := reflect.ValueOf(bindings)
	reflectValues = reflect.Indirect(reflectValues)
	if reflectValues.Kind() == reflect.Slice || reflectValues.Kind() == reflect.Array {
		for i := 0; i < reflectValues.Len(); i++ {
			value := reflectValues.Index(i).Interface()
			if !builder.isExpression(value) {
				values = append(values, value)
			}
		}
	}
	return values
}

// Get a scalar type value from an unknown type of input.
func (builder *Builder) flattenValue(value interface{}) interface{} {
	values := utils.Flatten(value)
	return values[0]
}

// MapScan scan the result from sql.Rows
func (builder *Builder) mapScan(rows *sql.Rows) ([]xun.R, error) {
	defer rows.Close()
	res := []xun.R{}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := builder.makeValues(len(columns))

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		dest := xun.R{}
		for i, column := range columns {
			dest[column] = builder.getValue(values[i])
		}
		res = append(res, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// structScan scan the result from sql.Rows
func (builder *Builder) structScan(rows *sql.Rows, v interface{}) error {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	structType, err := builder.getStructType(v)
	if err != nil {
		return err
	}

	fieldMap, err := builder.getFieldMap(structType)
	if err != nil {
		return err
	}

	vPtr := reflect.ValueOf(v)
	vRows := reflect.Indirect(vPtr)
	for rows.Next() {
		dest := reflect.New(structType)
		values, err := builder.makeScanValues(dest, fieldMap, columns)
		if err != nil {
			return err
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		value := reflect.Indirect(dest)
		vRows = reflect.Append(vRows, value)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	vPtr.Elem().Set(vRows)
	return nil
}

func (builder *Builder) getStructType(v interface{}) (reflect.Type, error) {

	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Thedest type is %s, it should be a pointer", reflectValue.Kind().String())
	}

	reflectValue = reflect.Indirect(reflectValue)
	structType := reflectValue.Type()
	if structType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Thedest type is %s, it should be a pointer of slice", reflectValue.Kind().String())
	}

	structType = reflectValue.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("non-struct dest type %s not support scan", reflectValue.Type().String())
	}
	return structType, nil
}

func (builder *Builder) getFieldMap(structType reflect.Type) (map[string]reflect.StructField, error) {
	fieldMap := map[string]reflect.StructField{}
	for i := 0; i < structType.NumField(); i++ {
		tag := xun.GetTagName(structType.Field(i), "json")
		if tag != "" && tag != "-" {
			fieldMap[tag] = structType.Field(i)
		}
	}
	return fieldMap, nil
}

func (builder *Builder) getValue(src interface{}) interface{} {
	value := src
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		value = reflect.Indirect(reflect.ValueOf(src)).Interface()
	}
	switch value.(type) {
	case []byte:
		return string(value.([]byte))
	default:
		return value
	}
}

func (builder *Builder) makeValues(length int) []interface{} {
	values := make([]interface{}, length)
	for i := range values {
		values[i] = new(interface{})
	}
	return values
}

func (builder *Builder) makeScanValues(dest reflect.Value, fieldMap map[string]reflect.StructField, columns []string) ([]interface{}, error) {

	values := []interface{}{}
	for _, column := range columns {
		field, has := fieldMap[column]
		if !has {
			return nil, fmt.Errorf("sql: expected %s destination arguments in Scan", column)
		}
		value := dest.Elem().FieldByName(field.Name)
		vPtr := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr()))
		values = append(values, vPtr.Interface())
	}

	return values, nil
}
