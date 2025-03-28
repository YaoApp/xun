package query

import "database/sql"

// Exec execute the sql, return the result
func (builder *Builder) Exec(sql string, bindings ...interface{}) (sql.Result, error) {
	stmt, err := builder.DB().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(bindings...)
}
