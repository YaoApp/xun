package query

import (
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Table create a new statement and set from givn table
func (builder *Builder) Table(name string) Query {
	builder.Query = dbal.NewQuery()
	builder.From(name)
	return builder
}

// Get Execute the query as a "select" statement.
func (builder *Builder) Get() ([]xun.R, error) {

	db := builder.DB()
	rows, err := db.Query(builder.ToSQL(), builder.GetBindings()...)
	if err != nil {
		return nil, err
	}

	return xun.MapScan(rows)

	// res := []xun.R{}
	// for rows.Next() {
	// 	row := map[string]interface{}{}
	// 	rows.MapScan(row)
	// 	utils.Println(row)
	// 	res = append(res, xun.MapToR(row))
	// }

	// return res, nil
}

// ToSQL Get the SQL representation of the query.
func (builder *Builder) ToSQL() string {
	return builder.Grammar.CompileSelect(builder.Query)
}

// GetBindings Get the current query value bindings in a flattened array.
func (builder *Builder) GetBindings() []interface{} {
	return builder.Query.GetBindings()
}

// MustGet Execute the query as a "select" statement.
func (builder *Builder) MustGet() []xun.R {
	res, err := builder.Get()
	utils.PanicIF(err)
	return res
}

// Find Execute a query for a single record by ID.
func (builder *Builder) Find() {
}

// MustFind  Execute a query for a single record by ID.
func (builder *Builder) MustFind() {
}

// Value Get a single column's value from the first result of a query.
func (builder *Builder) Value() {
}

// MustValue Get a single column's value from the first result of a query.
func (builder *Builder) MustValue() {
}

// Pluck Get an array with the values of a given column.
func (builder *Builder) Pluck() {
}

// MustPluck Get an array with the values of a given column.
func (builder *Builder) MustPluck() {
}

// Paginate paginate the given query into a simple paginator.
func (builder *Builder) Paginate() {
}

// MustPaginate paginate the given query into a simple paginator.
func (builder *Builder) MustPaginate() {
}

// When Executes the given closure when the first argument is true.
func (builder *Builder) When() {
}

// MustWhen Executes the given closure when the first argument is true.
func (builder *Builder) MustWhen() {
}

// Chunk Retrieves a small chunk of results at a time and feeds each chunk into a closure for processing.
func (builder *Builder) Chunk() {
}

// MustChunk Retrieves a small chunk of results at a time and feeds each chunk into a closure for processing.
func (builder *Builder) MustChunk() {
}

// Exists Determine if any rows exist for the current query.
func (builder *Builder) Exists() {
}

// MustExists Determine if any rows exist for the current query.
func (builder *Builder) MustExists() {
}

// DoesntExist Determine if no rows exist for the current query.
func (builder *Builder) DoesntExist() {
}

// MustDoesntExist Determine if no rows exist for the current query.
func (builder *Builder) MustDoesntExist() {
}
