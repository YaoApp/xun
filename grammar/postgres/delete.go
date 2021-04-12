package postgres

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// CompileDelete  Compile a delete statement into SQL.
func (grammarSQL Postgres) CompileDelete(query *dbal.Query) (string, []interface{}) {

	if len(query.Joins) == 0 && query.Limit < 0 {
		return grammarSQL.SQL.CompileDelete(query)
	}

	offset := 0
	bindings := []interface{}{}
	table := grammarSQL.WrapTable(query.From)

	alias := query.From.Alias
	if alias != "" {
		query.Columns = []interface{}{fmt.Sprintf("%s.ctid", alias)}
	} else {
		query.Columns = []interface{}{"ctid"}
	}

	selectSQL := grammarSQL.CompileSelectOffset(query, &offset)

	bindings = append(bindings, query.GetBindings()...)
	sql := fmt.Sprintf("delete from %s where %s in (%s)", table, grammarSQL.Wrap("ctid"), selectSQL)

	return sql, bindings
}
