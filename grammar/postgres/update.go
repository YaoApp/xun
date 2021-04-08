package postgres

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
)

// Upsert Upsert new records or update the existing ones.
func (grammarSQL Postgres) Upsert(query *dbal.Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (sql.Result, error) {
	sql, bindings := grammarSQL.CompileUpsert(query, values, uniqueBy, updateValues)
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())
	return grammarSQL.DB.Exec(sql, bindings...)
}

// CompileUpsert Upsert new records or update the existing ones.
func (grammarSQL Postgres) CompileUpsert(query *dbal.Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {

	if len(values) == 0 {
		return fmt.Sprintf("insert into %s default values", grammarSQL.WrapTable(query.From)), []interface{}{}
	}

	sql, bindings := grammarSQL.CompileInsert(query, values)
	sql = fmt.Sprintf("%s on conflict (%s) do update set", sql, grammarSQL.Columnize(uniqueBy))
	offset := len(bindings) + 1

	update := reflect.ValueOf(updateValues)
	kind := update.Kind()
	segments := []string{}
	if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < update.Len(); i++ {
			column := fmt.Sprintf("%v", update.Index(i).Interface())
			segments = append(segments, fmt.Sprintf("%s=excluded.%s", grammarSQL.Wrap(column), grammarSQL.Wrap(column)))
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
