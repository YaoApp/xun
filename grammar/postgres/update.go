package postgres

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
)

// Upsert Upsert new records or update the existing ones.
func (grammarSQL Postgres) Upsert(query *dbal.Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (sql.Result, error) {

	columns := values[0].Keys()
	insertValues := [][]interface{}{}
	for _, row := range values {
		insertValue := []interface{}{}
		for _, column := range columns {
			insertValue = append(insertValue, row.Get(column))
		}
		insertValues = append(insertValues, insertValue)
	}

	sql, bindings := grammarSQL.CompileUpsert(query, columns, insertValues, uniqueBy, updateValues)
	defer log.Debug(sql)
	return grammarSQL.DB.Exec(sql, bindings...)
}

// CompileUpsert Upsert new records or update the existing ones.
func (grammarSQL Postgres) CompileUpsert(query *dbal.Query, columns []interface{}, values [][]interface{}, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {

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
func (grammarSQL Postgres) CompileUpdate(query *dbal.Query, values map[string]interface{}) (string, []interface{}) {

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
		query.Columns = []interface{}{fmt.Sprintf("%s.ctid", alias)}
	} else {
		query.Columns = []interface{}{"ctid"}
	}

	selectSQL := grammarSQL.CompileSelectOffset(query, &offset)

	bindings = append(bindings, query.GetBindings()...)
	sql := fmt.Sprintf("update %s set %s where %s in (%s)", table, columns, grammarSQL.Wrap("ctid", true), selectSQL)

	return sql, bindings
}
