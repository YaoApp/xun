package query

import (
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// Update Update records in the database.
func (builder *Builder) Update(v interface{}) (int64, error) {

	values := xun.AnyToR(v).ToMap()
	sql, bindings := builder.Grammar.CompileUpdate(builder.Query, values)
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())

	res, err := builder.UseWrite().DB().Exec(sql, bindings...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// MustUpdate Update records in the database.
func (builder *Builder) MustUpdate(v interface{}) int64 {
	affected, err := builder.Update(v)
	utils.PanicIF(err)
	return affected
}

// UpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) UpdateOrInsert() {
}

// MustUpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) MustUpdateOrInsert() {
}

// Upsert new records or update the existing ones.
func (builder *Builder) Upsert(v interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) (int64, error) {

	columns, values := builder.prepareInsertValues(v, columns...)
	sql, bindings := builder.Grammar.CompileUpsert(builder.Query, columns, values, utils.Flatten(uniqueBy), update)
	defer logger.Debug(logger.UPDATE, sql).TimeCost(time.Now())

	res, err := builder.UseWrite().DB().Exec(sql, bindings...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// MustUpsert new records or update the existing ones.
func (builder *Builder) MustUpsert(values interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) int64 {
	affected, err := builder.Upsert(values, uniqueBy, update, columns...)
	utils.PanicIF(err)
	return affected
}

// Increment Increment a column's value by a given amount.
func (builder *Builder) Increment() {
}

// MustIncrement Increment a column's value by a given amount.
func (builder *Builder) MustIncrement() {
}

// Decrement Decrement a column's value by a given amount.
func (builder *Builder) Decrement() {
}

// MustDecrement Decrement a column's value by a given amount.
func (builder *Builder) MustDecrement() {
}
