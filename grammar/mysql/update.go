package mysql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileUpsert Upsert new records or update the existing ones.
func (grammarSQL MySQL) CompileUpsert(query *dbal.Query, columns []interface{}, values [][]interface{}, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {

	if len(values) == 0 {
		return fmt.Sprintf("insert into %s default values", grammarSQL.WrapTable(query.From)), []interface{}{}
	}

	sql, bindings := grammarSQL.CompileInsert(query, columns, values)
	sql = fmt.Sprintf("%s on duplicate key update", sql)
	offset := len(bindings)

	update := reflect.ValueOf(updateValues)
	kind := update.Kind()
	segments := []string{}
	if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < update.Len(); i++ {
			column := fmt.Sprintf("%v", update.Index(i).Interface())
			segments = append(segments, fmt.Sprintf("%s=values(%s)", grammarSQL.Wrap(column), grammarSQL.Wrap(column)))
		}
	} else if kind == reflect.Map {
		for _, key := range update.MapKeys() {
			column := fmt.Sprintf("%v", key)
			value := update.MapIndex(key).Interface()
			segments = append(segments, fmt.Sprintf("%s=%s", grammarSQL.Wrap(column), grammarSQL.Parameter(value, offset)))
			if !dbal.IsExpression(value) {
				bindings = append(bindings, value)
				offset++
			}
		}
	}

	return fmt.Sprintf("%s %s", sql, strings.Join(segments, ", ")), bindings
}
