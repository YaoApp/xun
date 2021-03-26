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
func (builder *Builder) InsertOrIgnore(v interface{}) (int64, error) {
	values := xun.AnyToRows(v)
	res, err := builder.Grammar.InsertIgnore(builder.Attr.From.TableFullName(), values)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, err
}

// MustInsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) MustInsertOrIgnore(v interface{}) int64 {
	affected, err := builder.InsertOrIgnore(v)
	utils.PanicIF(err)
	return affected
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
