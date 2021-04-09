package sqlite3

import (
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileInsertIgnore Compile an insert ignore statement into SQL.
func (grammarSQL SQLite3) CompileInsertIgnore(query *dbal.Query, columns []interface{}, values [][]interface{}) (string, []interface{}) {
	sql, bindings := grammarSQL.CompileInsert(query, columns, values)
	sql = strings.Replace(sql, "insert", "insert or ignore", 1)
	return sql, bindings
}
