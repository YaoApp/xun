package query

import (
	"reflect"
	"strings"
	"time"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/logger"
	"github.com/yaoapp/xun/utils"
)

// Insert Insert new records into the database.
func (builder *Builder) Insert(v interface{}) error {
	columns, values := builder.prepareInsertValues(v)
	sql, bindings := builder.Grammar.CompileInsert(builder.Query, columns, values)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())

	_, err := builder.UseWrite().DB().Exec(sql, bindings...)
	return err
}

// prepareInsertValues prepare the insert values
func (builder *Builder) prepareInsertValues(v interface{}, columns ...interface{}) ([]interface{}, [][]interface{}) {
	if _, ok := v.([][]interface{}); len(columns) > 1 && ok {
		return columns, v.([][]interface{})
	}
	values := xun.AnyToRows(v)
	columns = values[0].Keys()
	insertValues := [][]interface{}{}
	for _, row := range values {
		insertValue := []interface{}{}
		for _, column := range columns {
			insertValue = append(insertValue, row.MustGet(column))
		}
		insertValues = append(insertValues, insertValue)
	}
	return columns, insertValues
}

// MustInsert Insert new records into the database.
func (builder *Builder) MustInsert(v interface{}) {
	err := builder.Insert(v)
	utils.PanicIF(err)
}

// InsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) InsertOrIgnore(v interface{}) (int64, error) {
	columns, values := builder.prepareInsertValues(v)
	sql, bindings := builder.Grammar.CompileInsertIgnore(builder.Query, columns, values)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())

	res, err := builder.UseWrite().DB().Exec(sql, bindings...)
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
	seq := "id"
	if len(sequence) == 1 {
		seq = sequence[0]
	}
	columns, values := builder.prepareInsertValues(v)
	sql, bindings := builder.Grammar.CompileInsertGetID(builder.Query, columns, values, seq)
	defer logger.Debug(logger.CREATE, sql).TimeCost(time.Now())
	return builder.Grammar.ProcessInsertGetID(sql, bindings, seq)
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
	sql = builder.Grammar.CompileInsertUsing(builder.Query, columns, sql)
	res, err := builder.UseWrite().DB().Exec(sql, bindings...)
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
