package query

import (
	"fmt"
	"reflect"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Chunk Retrieves a small chunk of results at a time and feeds each chunk into a closure for processing.
func (builder *Builder) Chunk(size int, callback func(items []interface{}, page int) error, v ...interface{}) error {

	builder.enforceOrderBy()
	page := 1
	if size < 1 {
		size = 50
	}

	for {

		var results []interface{} = nil
		var countResults int

		// We'll execute the query for the given page and get the results. If there are
		// no results we can just break and return from here. When there are results
		// we will call the callback with the current chunk of these results here.
		if len(v) > 0 {

			reflectValuesPtr := reflect.ValueOf(v[0])
			reflectValues := reflect.Indirect(reflectValuesPtr)
			if reflectValues.Kind() != reflect.Slice {
				return fmt.Errorf("The given binding var shoule be a slice pointer")
			}

			reflectValuesType := reflectValues.Type()
			reflectValuesPtr.Elem().Set(reflect.New(reflectValuesType).Elem())

			_, err := builder.forPage(page, size).Get(v...)
			if err != nil {
				return err
			}

			countResults = reflectValues.Len()
			for i := 0; i < countResults; i++ {
				results = append(results, reflectValues.Index(i).Interface())
			}
		} else {
			rows, err := builder.forPage(page, size).Get()
			if err != nil {
				return err
			}

			countResults = len(rows)
			for _, row := range rows {
				results = append(results, row)
			}
		}

		// log.Trace("Chunk: countResults: %d size: %d page: %d", countResults, size, page)
		if err := callback(results, page); err != nil {
			// log.Error("Chunk: %s", err.Error())
			return err
		}

		if countResults == 0 || countResults != size {
			break
		}

		page++
	}

	return nil
}

// MustChunk Retrieves a small chunk of results at a time and feeds each chunk into a closure for processing.
func (builder *Builder) MustChunk(size int, callback func(items []interface{}, page int) error, v ...interface{}) {
	err := builder.Chunk(size, callback, v...)
	utils.PanicIF(err)
}

// ChunkByID chunk the results of a query by comparing IDs.
func (builder *Builder) ChunkByID() {}

// MustChunkByID chunk the results of a query by comparing IDs.
func (builder *Builder) MustChunkByID() {}

// Paginate paginate the given query into a simple paginator.
func (builder *Builder) Paginate(pageSize int, page int, v ...interface{}) (xun.P, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 15
	}

	total, err := builder.getCountForPagination([]interface{}{"*"})
	if err != nil {
		return xun.MakePaginator(0, pageSize, page), err
	}

	rows, err := builder.forPage(page, pageSize).Get(v...)
	if err != nil {
		return xun.MakePaginator(0, pageSize, page), err
	}

	items := []interface{}{}
	if rows != nil {
		for _, row := range rows {
			items = append(items, row)
		}
	} else if len(v) > 0 && reflect.TypeOf(v[0]).Kind() == reflect.Ptr {
		reflectRows := reflect.ValueOf(v[0])
		reflectRows = reflect.Indirect(reflectRows)
		if reflectRows.Kind() != reflect.Slice {
			return xun.MakeP(0, pageSize, page), fmt.Errorf("The given binding var shoule be a slice pointer")
		}
		for i := 0; i < reflectRows.Len(); i++ {
			items = append(items, reflectRows.Index(i).Interface())
		}
	}

	return xun.MakePaginator(total, pageSize, page, items...), nil
}

// MustPaginate paginate the given query into a simple paginator.
func (builder *Builder) MustPaginate(pageSize int, page int, v ...interface{}) xun.P {
	res, err := builder.Paginate(pageSize, page, v...)
	utils.PanicIF(err)
	return res
}

// Set the limit and offset for a given page.
func (builder *Builder) forPage(page int, pageSize int) Query {
	return builder.Offset((page - 1) * pageSize).Limit(pageSize)
}

// @todo: forPageBeforeID  Constrain the query to the previous "page" of results before a given ID.
// func (builder *Builder) forPageBeforeID(pageSize int, lastID int, column string) Query {
// 	builder.Query.Orders = builder.removeExistingOrdersFor(column)
// 	if lastID != 0 {
// 		builder.Where(column, "<", lastID)
// 	}
// 	return builder.OrderBy(column, "desc").Limit(pageSize)
// }

// @todo: forPageAfterID  Constrain the query to the next "page" of results after a given ID.
// func (builder *Builder) forPageAfterID(pageSize int, lastID int, column string) Query {
// 	builder.Query.Orders = builder.removeExistingOrdersFor(column)
// 	if lastID != 0 {
// 		builder.Where(column, ">", lastID)
// 	}
// 	return builder.OrderBy(column, "asc").Limit(pageSize)
// }

// getCountForPagination  Get the count of the total records for the paginator.
func (builder *Builder) getCountForPagination(columns []interface{}) (int, error) {

	if len(builder.Query.Groups) > 0 || len(builder.Query.Havings) > 0 {
		aggregate := 0
		clone := builder.cloneForPaginationCount()
		if len(clone.Query.Columns) == 0 && len(builder.Query.Joins) > 0 {
			if len(builder.Query.Groups) > 0 {
				clone.Select(builder.Query.Groups)
			} else if builder.Query.From.Alias != "" {
				clone.Select(dbal.Raw(fmt.Sprintf("%s.*", builder.Grammar.WrapTable(builder.Query.From.Alias))))
			} else {
				clone.Select(dbal.Raw(fmt.Sprintf("%s.*", builder.Grammar.WrapTable(builder.Query.From.Name))))
			}
		}

		_, err := builder.new().
			mergeBindings(clone).
			setAggregate("count", builder.withoutSelectAliases(columns)).
			FromRaw(fmt.Sprintf("(%s) as %s", clone.ToSQL(), builder.Grammar.Wrap("aggregate_table"))).
			Value("aggregate", &aggregate)

		return aggregate, err
	}

	clone := builder.cloneForPaginationCount()
	if len(builder.Query.Unions) == 0 {
		clone.Query.Columns = []interface{}{}
		clone.Query.Bindings["select"] = []interface{}{}
	}

	// fmt.Println(clone.setAggregate("count", builder.withoutSelectAliases(columns)).ToSQL())

	rows, err := clone.setAggregate("count", builder.withoutSelectAliases(columns)).Get()
	if err != nil {
		return 0, err
	}

	if len(rows) != 1 {
		return 0, nil
	}
	return int(rows[0].Get("aggregate").(int64)), nil
}

func (builder *Builder) cloneForPaginationCount() *Builder {
	new := builder.clone()
	new.Query.Orders = []dbal.Order{}
	new.Query.Limit = -1
	new.Query.Offset = -1
	new.Query.Bindings["order"] = []interface{}{}
	return new
}
