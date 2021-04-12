package sqlite3

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// CompileDelete Compile a delete statement into SQL.
func (grammarSQL SQLite3) CompileDelete(query *dbal.Query) (string, []interface{}) {

	if len(query.Joins) == 0 && query.Limit < 0 {
		return grammarSQL.SQL.CompileDelete(query)
	}

	offset := 0
	bindings := []interface{}{}
	table := grammarSQL.WrapTable(query.From)

	alias := query.From.Alias
	if alias != "" {
		query.Columns = []interface{}{fmt.Sprintf("%s.rowid", alias)}
	} else {
		query.Columns = []interface{}{"rowid"}
	}

	selectSQL := grammarSQL.CompileSelectOffset(query, &offset)

	bindings = append(bindings, query.GetBindings()...)
	sql := fmt.Sprintf("delete from  %s where %s in (%s)", table, grammarSQL.Wrap("rowid"), selectSQL)

	return sql, bindings
}

// CompileTruncate Compile a truncate table statement into SQL.
func (grammarSQL SQLite3) CompileTruncate(query *dbal.Query) ([]string, [][]interface{}) {

	return []string{
		fmt.Sprintf("delete from sqlite_sequence where name = ?"),
		fmt.Sprintf("delete from %s", grammarSQL.WrapTable(query.From)),
	}, [][]interface{}{{grammarSQL.getTableName(query)}, {}}
}

func (grammarSQL SQLite3) getTableName(query *dbal.Query) string {
	switch query.From.Name.(type) {
	case dbal.Name:
		return query.From.Name.(dbal.Name).Fullname()
	case dbal.Expression:
		return query.From.Name.(dbal.Expression).GetValue()
	case string:
		return query.From.Name.(string)
	default:
		return fmt.Sprintf("%v", query.From.Name)
	}
}
