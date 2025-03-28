package query

import "database/sql"

// Exec Use the current connection to execute the sql, return the result
func (builder *Builder) Exec(sql string, bindings ...interface{}) (sql.Result, error) {
	stmt, err := builder.DB().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(bindings...)
}

// ExecWrite Use the write connection to execute the sql, return the result
func (builder *Builder) ExecWrite(sql string, bindings ...interface{}) (sql.Result, error) {
	stmt, err := builder.DB(true).Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(bindings...)
}
