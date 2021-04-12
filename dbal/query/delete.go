package query

import (
	"time"

	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// Delete Delete records from the database.
func (builder *Builder) Delete() (int64, error) {
	sql, bindings := builder.Grammar.CompileDelete(builder.Query)
	defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())

	res, err := builder.UseWrite().DB().Exec(sql, bindings...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// MustDelete Delete records from the database.
func (builder *Builder) MustDelete() int64 {
	affected, err := builder.Delete()
	utils.PanicIF(err)
	return affected
}

// Truncate Run a truncate statement on the table.
func (builder *Builder) Truncate() error {
	sqls, bindings := builder.Grammar.CompileTruncate(builder.Query)
	for i, sql := range sqls {
		defer logger.Debug(logger.DELETE, sql).TimeCost(time.Now())
		_, err := builder.UseWrite().DB().Exec(sql, bindings[i]...)
		if err != nil {
			return err
		}
	}
	return nil
}

// MustTruncate Run a truncate statement on the table.
func (builder *Builder) MustTruncate() {
	err := builder.Truncate()
	utils.PanicIF(err)
}
