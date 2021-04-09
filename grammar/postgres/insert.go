package postgres

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// CompileInsertOrIgnore Compile an insert ignore statement into SQL.
func (grammarSQL Postgres) CompileInsertOrIgnore(query *dbal.Query, columns []interface{}, values [][]interface{}) (string, []interface{}) {
	sql, bindings := grammarSQL.CompileInsert(query, columns, values)
	sql = fmt.Sprintf("%s on conflict do nothing", sql)
	return sql, bindings
}

// CompileInsertGetID Compile an insert and get ID statement into SQL.
func (grammarSQL Postgres) CompileInsertGetID(query *dbal.Query, columns []interface{}, values [][]interface{}, sequence string) (string, []interface{}) {
	sql, bindings := grammarSQL.CompileInsert(query, columns, values)
	sql = fmt.Sprintf("%s returning %s", sql, grammarSQL.ID(sequence))
	return sql, bindings
}

// ProcessInsertGetID Execute an insert and get ID statement and return the id
func (grammarSQL Postgres) ProcessInsertGetID(sql string, bindings []interface{}, sequence string) (int64, error) {
	var seq int64
	err := grammarSQL.DB.Get(&seq, sql, bindings...)
	if err != nil {
		return 0, err
	}
	return seq, nil
}
