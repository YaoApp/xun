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
		safeFields = append(safeFields, grammarSQL.ID(field, grammarSQL.DB))
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName, grammarSQL.DB), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)

}

// InsertIgnore Insert ignore new records into the database.
func (grammarSQL SQL) InsertIgnore(tableName string, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field, grammarSQL.DB))
	}

	sql := fmt.Sprintf(`INSERT IGNORE INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName, grammarSQL.DB), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.NamedExec(sql, values)

}
