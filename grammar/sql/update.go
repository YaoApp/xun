package sql

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileUpsert Compile an "upsert" statement into SQL.
func (grammarSQL SQL) CompileUpsert(query *dbal.Query, columns []interface{}, values [][]interface{}, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {
	panic(fmt.Errorf("This database engine does not support upserts"))
}

// CompileUpdate Compile an update statement into SQL.
func (grammarSQL SQL) CompileUpdate(query *dbal.Query, values map[string]interface{}) (string, []interface{}) {

	offset := 0
	bindings := []interface{}{}
	table := grammarSQL.WrapTable(query.From)

	joins := ""
	if len(query.Joins) > 0 {
		joins = grammarSQL.CompileJoins(query, query.Joins, &offset)
		bindings = append(bindings, query.GetBindings("join")...)
		offset = len(bindings)
	}

	columns, columnsBindings := grammarSQL.CompileUpdateColumns(query, values, &offset)
	bindings = append(bindings, columnsBindings...)

	wheres := grammarSQL.CompileWheres(query, query.Wheres, &offset)
	bindings = append(bindings, query.GetBindings("where")...)

	return fmt.Sprintf("update %s %sset %s %s", table, joins, columns, wheres), bindings
}

// CompileUpdateColumns Compile the columns for an update statement.
func (grammarSQL SQL) CompileUpdateColumns(query *dbal.Query, values map[string]interface{}, offset *int) (string, []interface{}) {
	columns := []string{}
	bindings := []interface{}{}
	for key, value := range values {
		columns = append(columns, fmt.Sprintf("%s=%s", grammarSQL.Wrap(key, true), grammarSQL.Parameter(value, *offset+1)))
		if !dbal.IsExpression(value) {
			bindings = append(bindings, value)
			*offset++
		}
	}
	return strings.Join(columns, ", "), bindings
}
