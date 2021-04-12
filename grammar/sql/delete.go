package sql

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileDelete Compile a delete statement into SQL.
func (grammarSQL SQL) CompileDelete(query *dbal.Query) (string, []interface{}) {

	offset := 0
	bindings := []interface{}{}
	table := grammarSQL.WrapTable(query.From)
	alias := ""

	joins := ""
	if len(query.Joins) > 0 {
		joins = grammarSQL.CompileJoins(query, query.Joins, &offset)
		bindings = append(bindings, query.GetBindings("join")...)
		offset = len(bindings)
		tableArr := strings.Split(table, " as ")
		if len(tableArr) >= 1 {
			alias = tableArr[1]
		}
	}

	wheres := grammarSQL.CompileWheres(query, query.Wheres, &offset)
	bindings = append(bindings, query.GetBindings("where")...)

	if len(query.Joins) > 0 {
		return fmt.Sprintf("delete %s from %s %s %s", alias, table, joins, wheres), bindings
	}

	return fmt.Sprintf("delete from %s %s", table, wheres), bindings
}

// CompileTruncate Compile a truncate table statement into SQL.
func (grammarSQL SQL) CompileTruncate(query *dbal.Query) ([]string, [][]interface{}) {
	sql := fmt.Sprintf("truncate table %s", grammarSQL.WrapTable(query.From))
	return []string{sql}, [][]interface{}{{}}
}
