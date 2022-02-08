package query

import (
	"fmt"

	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Update Update records in the database.
func (builder *Builder) Update(v interface{}) (int64, error) {

	values := xun.MakeR(v).ToMap()
	sql, bindings := builder.Grammar.CompileUpdate(builder.Query, values)
	defer log.With(log.F{"bindings": bindings}).Debug(sql)

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

// MustUpdate Update records in the database.
func (builder *Builder) MustUpdate(v interface{}) int64 {
	affected, err := builder.Update(v)
	utils.PanicIF(err)
	return affected
}

// UpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) UpdateOrInsert(attributes interface{}, values ...interface{}) (bool, error) {

	exists, err := builder.Where(attributes).Exists()
	if err != nil {
		return false, err
	}

	if !exists {
		insertValues := xun.MakeR(attributes)
		if len(values) > 0 {
			insertValues.Merge(values[0])
		}
		err = builder.Insert(insertValues)
		if err != nil {
			return false, err
		}
	}

	if len(values) == 0 {
		return true, nil
	}

	updateValues := values[0]
	_, err = builder.Limit(1).Update(updateValues)
	if err != nil {
		return false, err
	}
	return true, nil
}

// MustUpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) MustUpdateOrInsert(attributes interface{}, values ...interface{}) bool {
	res, err := builder.UpdateOrInsert(attributes, values...)
	utils.PanicIF(err)
	return res
}

// Upsert new records or update the existing ones.
func (builder *Builder) Upsert(v interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) (int64, error) {

	columns, values := builder.prepareInsertValues(v, columns...)
	sql, bindings := builder.Grammar.CompileUpsert(builder.Query, columns, values, utils.Flatten(uniqueBy), update)
	defer log.With(log.F{"bindings": bindings}).Debug(sql)

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

// MustUpsert new records or update the existing ones.
func (builder *Builder) MustUpsert(values interface{}, uniqueBy interface{}, update interface{}, columns ...interface{}) int64 {
	affected, err := builder.Upsert(values, uniqueBy, update, columns...)
	utils.PanicIF(err)
	return affected
}

// Increment Increment a column's value by a given amount.
func (builder *Builder) Increment(column interface{}, amount interface{}, extra ...interface{}) (int64, error) {
	if !utils.IsNumeric(amount) {
		panic(fmt.Errorf("non-numeric value passed to increment method"))
	}
	wrapped := builder.Grammar.Wrap(column)
	values := map[string]interface{}{}
	if len(extra) > 0 {
		values = xun.MakeR(extra[0]).ToMap()
	}
	values[wrapped] = dbal.Raw(fmt.Sprintf("%s+%v", wrapped, amount))
	return builder.Update(values)
}

// MustIncrement Increment a column's value by a given amount.
func (builder *Builder) MustIncrement(column interface{}, amount interface{}, extra ...interface{}) int64 {
	affected, err := builder.Increment(column, amount, extra...)
	utils.PanicIF(err)
	return affected
}

// Decrement Decrement a column's value by a given amount.
func (builder *Builder) Decrement(column interface{}, amount interface{}, extra ...interface{}) (int64, error) {
	if !utils.IsNumeric(amount) {
		panic(fmt.Errorf("non-numeric value passed to decrement method"))
	}
	wrapped := builder.Grammar.Wrap(column)
	values := map[string]interface{}{}
	if len(extra) > 0 {
		values = xun.MakeR(extra[0]).ToMap()
	}
	values[wrapped] = dbal.Raw(fmt.Sprintf("%s-%v", wrapped, amount))
	return builder.Update(values)
}

// MustDecrement Decrement a column's value by a given amount.
func (builder *Builder) MustDecrement(column interface{}, amount interface{}, extra ...interface{}) int64 {
	affected, err := builder.Decrement(column, amount, extra...)
	utils.PanicIF(err)
	return affected
}
