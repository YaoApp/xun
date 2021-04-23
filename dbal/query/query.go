package query

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// Table create a new statement and set from givn table
func (builder *Builder) Table(name string) Query {
	builder.Query = dbal.NewQuery()
	builder.From(name)
	return builder
}

// Get Execute the query as a "select" statement.
func (builder *Builder) Get(v ...interface{}) ([]xun.R, error) {
	db := builder.DB()
	stmt, err := db.Prepare(builder.ToSQL())
	if err != nil {
		defer logger.Debug(logger.RETRIEVE, builder.ToSQL(), fmt.Sprintf("%v", builder.GetBindings())).Write()
		defer logger.Fatal(500, err.Error()).Write()
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(builder.GetBindings()...)
	if err != nil {
		return nil, err
	}

	if len(v) == 1 && v[0] != nil {
		if reflect.TypeOf(v[0]).Kind() != reflect.Ptr {
			return nil, fmt.Errorf("The input param is %s, it should be a pointer", reflect.TypeOf(v[0]).Kind().String())
		}
		err := builder.structScan(rows, v[0])
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return builder.mapScan(rows)
}

// MustGet Execute the query as a "select" statement.
func (builder *Builder) MustGet(v ...interface{}) []xun.R {
	res, err := builder.Get(v...)
	utils.PanicIF(err)
	return res
}

// First Execute the query and get the first result.
func (builder *Builder) First(v ...interface{}) (xun.R, error) {
	rows, err := builder.Take(1).Get(v...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 1 {
		return rows[0], err
	}
	return xun.MakeR(), nil
}

// MustFirst Execute the query and get the first result.
func (builder *Builder) MustFirst(v ...interface{}) xun.R {
	res, err := builder.First(v...)
	utils.PanicIF(err)
	return res
}

// ToSQL Get the SQL representation of the query.
func (builder *Builder) ToSQL() string {
	return builder.Grammar.CompileSelect(builder.Query)
}

// GetBindings Get the current query value bindings in a flattened array.
func (builder *Builder) GetBindings() []interface{} {
	return builder.Query.GetBindings()
}

// Exists Determine if any rows exist for the current query.
func (builder *Builder) Exists() (bool, error) {
	sql := builder.Grammar.CompileExists(builder.Query)

	db := builder.DB()
	rows, err := db.Query(sql, builder.GetBindings()...)
	if err != nil {
		return false, err
	}

	res, err := builder.mapScan(rows)
	if err != nil {
		return false, err
	}

	if len(res) == 1 {
		exists := fmt.Sprintf("%v", res[0]["exists"])
		return (exists == "1" || exists == "true"), nil
	}

	return false, nil
}

// MustExists Determine if any rows exist for the current query.
func (builder *Builder) MustExists() bool {
	res, err := builder.Exists()
	utils.PanicIF(err)
	return res
}

// DoesntExist Determine if no rows exist for the current query.
func (builder *Builder) DoesntExist() (bool, error) {
	res, err := builder.Exists()
	if err != nil {
		return false, err
	}
	return !res, nil
}

// MustDoesntExist Determine if no rows exist for the current query.
func (builder *Builder) MustDoesntExist() bool {
	res, err := builder.DoesntExist()
	utils.PanicIF(err)
	return res
}

// Find Execute a query for a single record by ID.
func (builder *Builder) Find(id interface{}, args ...interface{}) (xun.R, error) {
	key, v := builder.prepareFindArgs(args...)
	return builder.Where(key, id).First(v)
}

// MustFind  Execute a query for a single record by ID.
func (builder *Builder) MustFind(id interface{}, args ...interface{}) xun.R {
	res, err := builder.Find(id, args...)
	utils.PanicIF(err)
	return res
}

// Value Get a single column's value from the first result of a query.
func (builder *Builder) Value(column string, v ...interface{}) (interface{}, error) {
	row, err := builder.Select(column).First(v...)
	if err != nil {
		return nil, err
	}

	if row == nil || row.IsEmpty() {
		return nil, nil
	}

	return row.Get(column), nil
}

// MustValue Get a single column's value from the first result of a query.
func (builder *Builder) MustValue(column string, v ...interface{}) interface{} {
	res, err := builder.Value(column, v...)
	utils.PanicIF(err)
	return res
}

// Pluck Get an array with the values of a given column.
func (builder *Builder) Pluck() {
}

// MustPluck Get an array with the values of a given column.
func (builder *Builder) MustPluck() {
}
