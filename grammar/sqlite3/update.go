package sqlite3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileUpsert Upsert new records or update the existing ones.
func (grammarSQL SQLite3) CompileUpsert(query *dbal.Query, columns []interface{}, values [][]interface{}, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {

	if len(values) == 0 {
		return fmt.Sprintf("insert into %s default values", grammarSQL.WrapTable(query.From)), []interface{}{}
	}

	sql, bindings := grammarSQL.CompileInsert(query, columns, values)
	sql = fmt.Sprintf("%s on conflict (%s) do update set", sql, grammarSQL.Columnize(uniqueBy))
	offset := len(bindings) + 1

	update := reflect.ValueOf(updateValues)
	kind := update.Kind()
	segments := []string{}
	if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < update.Len(); i++ {
			column := fmt.Sprintf("%v", update.Index(i).Interface())
			segments = append(segments, fmt.Sprintf("%s=excluded.%s", grammarSQL.Wrap(column, true), grammarSQL.Wrap(column, true)))
		}
	} else if kind == reflect.Map {
		for _, key := range update.MapKeys() {
			column := fmt.Sprintf("%v", key)
			value := update.MapIndex(key).Interface()
			segments = append(segments, fmt.Sprintf("%s=%s", grammarSQL.Wrap(column, true), grammarSQL.Parameter(value, offset)))
			if !dbal.IsExpression(value) {
				bindings = append(bindings, value)
				offset++
			}
		}
	}

	return fmt.Sprintf("%s %s", sql, strings.Join(segments, ", ")), bindings
}

// CompileUpdate Compile an update statement into SQL.
func (grammarSQL SQLite3) CompileUpdate(query *dbal.Query, values map[string]interface{}) (string, []interface{}) {

	if len(query.Joins) == 0 && query.Limit < 0 {
		return grammarSQL.SQL.CompileUpdate(query, values)
	}

	offset := 0
	bindings := []interface{}{}
	table := grammarSQL.WrapTable(query.From)

	columns, columnsBindings := grammarSQL.CompileUpdateColumns(query, values, &offset)
	bindings = append(bindings, columnsBindings...)

	alias := query.From.Alias
	if alias != "" {
		query.Columns = []interface{}{fmt.Sprintf("%s.rowid", alias)}
	} else {
		query.Columns = []interface{}{"rowid"}
	}

	selectSQL := grammarSQL.CompileSelectOffset(query, &offset)

	bindings = append(bindings, query.GetBindings()...)
	sql := fmt.Sprintf("update %s set %s where %s in (%s)", table, columns, grammarSQL.Wrap("rowid", true), selectSQL)

	return sql, bindings
}
