package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/logger"
)

// Insert Insert new records into the database.
func (grammarSQL SQL) Insert(tableName string, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field))
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)
}

// InsertIgnore Insert ignore new records into the database.
func (grammarSQL SQL) InsertIgnore(tableName string, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field))
	}

	sql := fmt.Sprintf(`INSERT IGNORE INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)

}

// InsertGetID Insert new records into the database and return the last insert ID
func (grammarSQL SQL) InsertGetID(tableName string, values []xun.R, sequence string) (int64, error) {
	res, err := grammarSQL.Insert(tableName, values)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
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
