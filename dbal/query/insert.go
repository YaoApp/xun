package query

import (
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/utils"
)

// Insert Insert new records into the database.
func (builder *Builder) Insert(v interface{}) error {
	values := xun.AnyToRows(v)
	_, err := builder.Grammar.Insert(builder.Attr.From.TableFullName(), values)
	return err
}

// MustInsert Insert new records into the database.
func (builder *Builder) MustInsert(v interface{}) {
	err := builder.Insert(v)
	utils.PanicIF(err)
}

// InsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) InsertOrIgnore() {
}

// MustInsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) MustInsertOrIgnore() {
}

// InsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) InsertGetID() {
}

// MustInsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) MustInsertGetID() {
}

// InsertUsing Insert new records into the table using a subquery.
func (builder *Builder) InsertUsing() {
}

// MustInsertUsing Insert new records into the table using a subquery.
func (builder *Builder) MustInsertUsing() {
}
