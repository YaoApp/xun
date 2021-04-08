package query

import (
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/utils"
)

// Update Update records in the database.
func (builder *Builder) Update() {
}

// MustUpdate Update records in the database.
func (builder *Builder) MustUpdate() {
}

// UpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) UpdateOrInsert() {
}

// MustUpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) MustUpdateOrInsert() {
}

// Upsert new records or update the existing ones.
func (builder *Builder) Upsert(values interface{}, uniqueBy interface{}, update interface{}) (int64, error) {
	builder.UseWrite()
	res, err := builder.Grammar.Upsert(builder.Query, xun.AnyToRows(values), utils.Flatten(uniqueBy), update)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// MustUpsert new records or update the existing ones.
func (builder *Builder) MustUpsert(values interface{}, uniqueBy interface{}, update interface{}) int64 {
	affected, err := builder.Upsert(values, uniqueBy, update)
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
