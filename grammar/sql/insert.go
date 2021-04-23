package sql

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
)

// CompileInsert Compile an insert statement into SQL.
func (grammarSQL SQL) CompileInsert(query *dbal.Query, columns []interface{}, values [][]interface{}) (string, []interface{}) {

	table := grammarSQL.WrapTable(query.From)
	if len(values) == 0 {
		return fmt.Sprintf("insert into %s default values", table), nil
	}

	offset := 0
	parameters := []string{}
	bindings := []interface{}{}
	for _, value := range values {
		parameters = append(parameters, fmt.Sprintf("(%s)", grammarSQL.Parameterize(value, offset)))
		for _, v := range value {
			if !dbal.IsExpression(v) {
				bindings = append(bindings, v)
				offset++
			}
		}
	}

	return fmt.Sprintf("insert into %s (%s) values %s", table, grammarSQL.Columnize(columns), strings.Join(parameters, ",")), bindings
}

// CompileInsertOrIgnore Compile an insert ignore statement into SQL.
func (grammarSQL SQL) CompileInsertOrIgnore(query *dbal.Query, columns []interface{}, values [][]interface{}) (string, []interface{}) {
	panic(fmt.Errorf("This database engine does not support upserts"))
}

// CompileInsertGetID Compile an insert and get ID statement into SQL.
func (grammarSQL SQL) CompileInsertGetID(query *dbal.Query, columns []interface{}, values [][]interface{}, sequence string) (string, []interface{}) {
	return grammarSQL.CompileInsert(query, columns, values)
}

// ProcessInsertGetID Execute an insert and get ID statement and return the id
func (grammarSQL SQL) ProcessInsertGetID(sql string, bindings []interface{}, sequence string) (int64, error) {
	stmt, err := grammarSQL.DB.Prepare(sql)
	defer stmt.Close()

	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(bindings...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// CompileInsertUsing Compile an insert statement using a subquery into SQL.
func (grammarSQL SQL) CompileInsertUsing(query *dbal.Query, columns []interface{}, sql string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) %s", grammarSQL.WrapTable(query.From), grammarSQL.Columnize(columns), sql)
}
