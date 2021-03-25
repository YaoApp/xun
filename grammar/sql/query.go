package sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/logger"
)

// Insert Insert new records into the database.
func (grammarSQL SQL) Insert(tableName string, v interface{}) (sql.Result, error) {

	values := []xun.R{}
	switch v.(type) {
	case xun.R:
		value := v.(xun.R)
		values = append(values, value)
		break
	case []xun.R:
		values = v.([]xun.R)
		if len(values) < 1 {
			return nil, fmt.Errorf("Insert into %s error. The input values is empty", tableName)
		}
		break
	default:
		var err error
		var value xun.R
		kind := reflect.TypeOf(v).Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			values, err = xun.AnyToRs(v)
		} else {
			value, err = xun.AnyToR(v)
			values = append(values, value)
		}
		if err != nil {
			return nil, fmt.Errorf("Insert into %s error. %s", tableName, err)
		}
	}

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field, grammarSQL.DB))
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName, grammarSQL.DB), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)

}
