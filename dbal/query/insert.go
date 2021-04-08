package query

import (
	"reflect"
	"strings"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/utils"
)

// Insert Insert new records into the database.
func (builder *Builder) Insert(v interface{}) error {
	builder.UseWrite()
	values := xun.AnyToRows(v)
	_, err := builder.Grammar.Insert(builder.Query, values)
	return err
}

// MustInsert Insert new records into the database.
func (builder *Builder) MustInsert(v interface{}) {
	err := builder.Insert(v)
	utils.PanicIF(err)
}

// InsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) InsertOrIgnore(v interface{}) (int64, error) {
	builder.UseWrite()
	values := xun.AnyToRows(v)
	res, err := builder.Grammar.InsertIgnore(builder.Query, values)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// MustInsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) MustInsertOrIgnore(v interface{}) int64 {
	affected, err := builder.InsertOrIgnore(v)
	utils.PanicIF(err)
	return affected
}

// InsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) InsertGetID(v interface{}, sequence ...string) (int64, error) {
	builder.UseWrite()
	values := xun.AnyToRows(v)
	seq := "id"
	if len(sequence) == 1 {
		seq = sequence[0]
	}
	return builder.Grammar.InsertGetID(builder.Query, values, seq)
}

// MustInsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) MustInsertGetID(v interface{}, sequence ...string) int64 {
	lastID, err := builder.InsertGetID(v, sequence...)
	utils.PanicIF(err)
	return lastID
}

// InsertUsing Insert new records into the table using a subquery.
func (builder *Builder) InsertUsing(qb interface{}, columns ...interface{}) (int64, error) {

	// columns  "field1,field2", []string{"field1", "field2"}
	if len(columns) == 1 {
		col, ok := columns[0].(string)
		if ok && strings.Contains(col, ",") {
			cols := strings.Split(col, ",")
			columns = []interface{}{}
			for _, col := range cols {
				columns = append(columns, strings.Trim(col, " "))
			}
		} else if !ok {
			reflectValue := reflect.ValueOf(columns[0])
			kind := reflectValue.Kind()
			columns = []interface{}{}
			if kind == reflect.Array || kind == reflect.Slice {
				for i := 0; i < reflectValue.Len(); i++ {
					columns = append(columns, reflectValue.Index(i).Interface())
				}
			}

		}
	}

	builder.UseWrite()
	sql, bindings, _ := builder.createSub(qb)
	res, err := builder.Grammar.InsertUsing(builder.Query, columns, sql, bindings)
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
