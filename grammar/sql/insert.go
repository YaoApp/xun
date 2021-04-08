package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
)

// Insert Insert new records into the database.
func (grammarSQL SQL) Insert(query *dbal.Query, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field))
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, grammarSQL.WrapTable(query.From), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)
}

// CompileInsert Compile an insert statement into SQL.
func (grammarSQL SQL) CompileInsert(query *dbal.Query, values []xun.R) (string, []interface{}) {

	table := grammarSQL.WrapTable(query.From)
	if len(values) == 0 {
		return fmt.Sprintf("insert into %s default values", table), nil
	}

	columns := values[0].Keys()
	insertValues := [][]interface{}{}
	for _, row := range values {
		insertValue := []interface{}{}
		for _, column := range columns {
			insertValue = append(insertValue, row.MustGet(column))
		}
		insertValues = append(insertValues, insertValue)
	}

	offset := 0
	parameters := []string{}
	bindings := []interface{}{}
	for _, value := range insertValues {
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

// InsertIgnore Insert ignore new records into the database.
func (grammarSQL SQL) InsertIgnore(query *dbal.Query, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field))
	}

	sql := fmt.Sprintf(`INSERT IGNORE INTO %s (%s) VALUES (%s)`, grammarSQL.WrapTable(query.From), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)

}

// InsertGetID Insert new records into the database and return the last insert ID
func (grammarSQL SQL) InsertGetID(query *dbal.Query, values []xun.R, sequence string) (int64, error) {
	res, err := grammarSQL.Insert(query, values)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// InsertUsing Compile and run an insert statement using a subquery into SQL.
func (grammarSQL SQL) InsertUsing(query *dbal.Query, columns []interface{}, sql string, bindings []interface{}) (sql.Result, error) {
	sql = fmt.Sprintf("INSERT INTO %s (%s) %s", grammarSQL.WrapTable(query.From), grammarSQL.Columnize(columns), sql)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	return grammarSQL.DB.Exec(sql, bindings...)
}

// GetOperators get the operators
func (grammarSQL SQL) GetOperators() []string {
	return []string{
		"=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
		"like", "like binary", "not like", "ilike",
		"&", "|", "^", "<<", ">>",
		"rlike", "not rlike", "regexp", "not regexp",
		"~", "~*", "!~", "!~*", "similar to",
		"not similar to", "not ilike", "~~*", "!~~*",
	}
}

//GetSelectComponents Get The components that make up a select clause.
func (grammarSQL SQL) GetSelectComponents() []string {
	return []string{"aggregate", "columns", "from", "joins", "wheres", "groups", "havings", "orders", "limit", "offset", "lock"}
}
