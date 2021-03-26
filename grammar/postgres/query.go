package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/logger"
)

// InsertIgnore Insert ignore new records into the database.
func (grammarSQL Postgres) InsertIgnore(tableName string, values []xun.R) (sql.Result, error) {

	safeFields := []string{}
	bindVars := []string{}
	for field := range values[0] {
		bindVars = append(bindVars, ":"+field)
		safeFields = append(safeFields, grammarSQL.ID(field, grammarSQL.DB))
	}
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, grammarSQL.ID(tableName, grammarSQL.DB), strings.Join(safeFields, ","), strings.Join(bindVars, ","))
	sql, args, _ := grammarSQL.DB.BindNamed(sql, values)
	sql = sql + " ON CONFLICT DO NOTHING"
	defer logger.Debug(logger.RETRIEVE, sql).TimeCost(time.Now())
	return grammarSQL.DB.Exec(sql, args...)

}
