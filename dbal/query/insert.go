package query

import (
	"fmt"
	"time"

	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// Insert Insert new records into the database.
func (builder *Builder) Insert(v interface{}, columns ...interface{}) error {
	columns, values := builder.prepareInsertValues(v, columns...)
	sql, bindings := builder.Grammar.CompileInsert(builder.Query, columns, values)
	defer logger.Debug(logger.CREATE, sql, fmt.Sprintf("%v", bindings)).TimeCost(time.Now())

	stmt, err := builder.UseWrite().DB().Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bindings...)
	return err
}

// MustInsert Insert new records into the database.
func (builder *Builder) MustInsert(v interface{}, columns ...interface{}) {
	err := builder.Insert(v, columns...)
	utils.PanicIF(err)
}

// InsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) InsertOrIgnore(v interface{}, columns ...interface{}) (int64, error) {
	columns, values := builder.prepareInsertValues(v, columns...)
	sql, bindings := builder.Grammar.CompileInsertOrIgnore(builder.Query, columns, values)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())

	stmt, err := builder.UseWrite().DB().Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bindings...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// MustInsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) MustInsertOrIgnore(v interface{}, columns ...interface{}) int64 {
	affected, err := builder.InsertOrIgnore(v, columns...)
	utils.PanicIF(err)
	return affected
}

// InsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) InsertGetID(v interface{}, args ...interface{}) (int64, error) {
	seq := "id"
	columns := []interface{}{}

	if len(args) == 1 {
		seq = fmt.Sprintf("%v", args[0])
	} else if len(args) > 1 {
		columns = args[1:]
	}

	columns, values := builder.prepareInsertValues(v, columns...)
	sql, bindings := builder.Grammar.CompileInsertGetID(builder.Query, columns, values, seq)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	return builder.Grammar.ProcessInsertGetID(sql, bindings, seq)
}

// MustInsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) MustInsertGetID(v interface{}, args ...interface{}) int64 {
	lastID, err := builder.InsertGetID(v, args...)
	utils.PanicIF(err)
	return lastID
}

// InsertUsing Insert new records into the table using a subquery.
func (builder *Builder) InsertUsing(qb interface{}, columns ...interface{}) (int64, error) {

	columns = builder.prepareColumns(columns...)
	sub, bindings, _ := builder.createSub(qb)
	sql := builder.parseSub(sub)
	sql = builder.Grammar.CompileInsertUsing(builder.Query, columns, sql)

	stmt, err := builder.UseWrite().DB().Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bindings...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// MustInsertUsing Insert new records into the table using a subquery.
func (builder *Builder) MustInsertUsing(qb interface{}, columns ...interface{}) int64 {
	affected, err := builder.InsertUsing(qb, columns...)
	utils.PanicIF(err)
	return affected
}
