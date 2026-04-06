package sqlite3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	goSQL "github.com/yaoapp/xun/grammar/sql"
)

func newTestSQLite3() SQLite3 {
	return SQLite3{
		SQL: goSQL.NewSQL(&goSQL.Quoter{}),
	}
}

func newTestSQLite3WithQuoter() SQLite3 {
	return SQLite3{
		SQL: goSQL.NewSQL(&Quoter{}),
	}
}

func newBaseQuery(tableName string) *dbal.Query {
	return &dbal.Query{
		From:   dbal.From{Type: "table", Name: dbal.NewName(tableName)},
		Limit:  -1,
		Offset: -1,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {}, "groupBy": {}, "having": {}, "order": {},
		},
	}
}

func TestWhereJsoncontainsSQLite(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%admin%",
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "`tags` like ?", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsSQLiteNot(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%admin%",
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "not `tags` like ?", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsSQLiteOffset(t *testing.T) {
	g := newTestSQLite3()
	offset := 5
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%test%",
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "`tags` like ?", result)
	assert.Equal(t, 6, offset)
}

func newDeleteQuery(tableName string, wheres []dbal.Where, whereBindings []interface{}) *dbal.Query {
	return &dbal.Query{
		From:   dbal.From{Name: dbal.NewName(tableName)},
		Limit:  -1,
		Wheres: wheres,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where":   whereBindings,
			"groupBy": {}, "having": {}, "order": {},
		},
	}
}

func TestCompileDeleteSimpleSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("users", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := g.CompileDelete(query)
	assert.Equal(t, "delete from `users` where `id` = ?", sql)
	assert.Equal(t, []interface{}{1}, bindings)
}

func TestCompileDeleteWithJsonContainsSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("assistants", []dbal.Where{
		{Type: "jsoncontains", Column: "tags", Value: "%test%", Boolean: "and", Offset: 1},
	}, []interface{}{"%test%"})

	sql, bindings := g.CompileDelete(query)
	assert.Equal(t, "delete from `assistants` where `tags` like ?", sql)
	assert.Equal(t, []interface{}{"%test%"}, bindings)
}

func TestCompileDeleteNoWhereSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("users", nil, nil)
	query.Bindings["where"] = []interface{}{}

	sql, bindings := g.CompileDelete(query)
	assert.Equal(t, "delete from `users` ", sql)
	assert.Empty(t, bindings)
}

func TestCompileUpdateSimpleSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("users", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := g.CompileUpdate(query, map[string]interface{}{"name": "new_name"})
	assert.Contains(t, sql, "update `users` set")
	assert.Contains(t, sql, "where `id` =")
	assert.Contains(t, bindings, 1)
	assert.Contains(t, bindings, "new_name")
}

func TestCompileUpdateWithJsonContainsSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("assistants", []dbal.Where{
		{Type: "jsoncontains", Column: "tags", Value: "%old%", Boolean: "and", Offset: 1},
	}, []interface{}{"%old%"})

	sql, bindings := g.CompileUpdate(query, map[string]interface{}{"name": "updated"})
	assert.Contains(t, sql, "update `assistants` set")
	assert.Contains(t, sql, "`tags` like")
	assert.Contains(t, bindings, "%old%")
	assert.Contains(t, bindings, "updated")
}

func newQueryWithLimit(tableName string, limit int, wheres []dbal.Where, whereBindings []interface{}) *dbal.Query {
	return &dbal.Query{
		From:   dbal.From{Name: dbal.NewName(tableName)},
		Limit:  limit,
		Wheres: wheres,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where":   whereBindings,
			"groupBy": {}, "having": {}, "order": {},
		},
	}
}

func newQueryWithAlias(tableName, alias string, wheres []dbal.Where, whereBindings []interface{}) *dbal.Query {
	return &dbal.Query{
		From:   dbal.From{Name: dbal.NewName(tableName), Alias: alias},
		Limit:  10,
		Wheres: wheres,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where":   whereBindings,
			"groupBy": {}, "having": {}, "order": {},
		},
	}
}

func TestCompileDeleteWithLimitSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newQueryWithLimit("users", 10, []dbal.Where{
		{Type: "basic", Column: "status", Operator: "=", Value: "inactive", Boolean: "and", Offset: 1},
	}, []interface{}{"inactive"})

	sql, bindings := g.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, "`rowid`")
	assert.Contains(t, sql, "in (")
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteWithLimitAliasSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newQueryWithAlias("users", "u", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := g.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, "`rowid`")
	assert.Contains(t, sql, "`u`.`rowid`")
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteWithLimitNoWhereSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newQueryWithLimit("users", 5, nil, nil)
	query.Bindings["where"] = []interface{}{}

	sql, _ := g.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, "`rowid`")
}

func TestCompileUpdateWithLimitSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newQueryWithLimit("users", 10, []dbal.Where{
		{Type: "basic", Column: "status", Operator: "=", Value: "old", Boolean: "and", Offset: 1},
	}, []interface{}{"old"})

	sql, bindings := g.CompileUpdate(query, map[string]interface{}{"status": "new"})
	assert.Contains(t, sql, "update")
	assert.Contains(t, sql, "`rowid`")
	assert.Contains(t, sql, "in (")
	assert.NotEmpty(t, bindings)
}

func TestCompileUpdateWithLimitAliasSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newQueryWithAlias("users", "u", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := g.CompileUpdate(query, map[string]interface{}{"name": "updated"})
	assert.Contains(t, sql, "update")
	assert.Contains(t, sql, "`rowid`")
	assert.Contains(t, sql, "`u`.`rowid`")
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteNestedJsonContainsSQLite(t *testing.T) {
	g := newTestSQLite3()
	query := newDeleteQuery("assistants", []dbal.Where{
		{
			Type:    "nested",
			Boolean: "and",
			Query: &dbal.Query{
				Wheres: []dbal.Where{
					{Type: "jsoncontains", Column: "tags", Value: "%a%", Boolean: "and", Offset: 1},
					{Type: "jsoncontains", Column: "tags", Value: "%b%", Boolean: "or", Offset: 1},
				},
				Bindings: map[string][]interface{}{
					"select": {}, "from": {}, "join": {},
					"where":   {"%a%", "%b%"},
					"groupBy": {}, "having": {}, "order": {},
				},
			},
		},
	}, []interface{}{"%a%", "%b%"})

	sql, bindings := g.CompileDelete(query)
	assert.Contains(t, sql, "`tags` like")
	assert.Contains(t, sql, "delete from")
	assert.Len(t, bindings, 2)
}

// ---------------------------------------------------------------------------
// CompileSelect / CompileSelectOffset
// ---------------------------------------------------------------------------

func TestCompileSelectBasic(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	result := g.CompileSelect(query)
	assert.Contains(t, result, "select")
	assert.Contains(t, result, "from `users`")
}

func TestCompileSelectWithSQL(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.SQL = "select id, name from users where id = 1"
	result := g.CompileSelect(query)
	assert.Equal(t, "select id, name from users where id = 1", result)
}

func TestCompileSelectWithSQLContainingLimit(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.SQL = "select * from users limit 10"
	result := g.CompileSelect(query)
	assert.Equal(t, "select * from users limit 10", result)
}

func TestCompileSelectWithSQLContainingOffset(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.SQL = "select * from users offset 5"
	result := g.CompileSelect(query)
	assert.Equal(t, "select * from users offset 5", result)
}

func TestCompileSelectWithSQLAppendsLimitOffset(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.SQL = "select * from users where active = 1"
	query.Limit = 10
	query.Offset = 20
	result := g.CompileSelect(query)
	assert.Contains(t, result, "select * from users where active = 1")
	assert.Contains(t, result, "limit")
}

func TestCompileSelectUnionAggregate(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Aggregate = dbal.Aggregate{Func: "count", Columns: []interface{}{"*"}}
	unionQuery := newBaseQuery("admins")
	query.Unions = []dbal.Union{{All: false, Query: unionQuery}}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "count")
	assert.Contains(t, result, "temp_table")
}

func TestCompileSelectWithUnions(t *testing.T) {
	g := newTestSQLite3WithQuoter()
	query := newBaseQuery("users")
	unionQuery := newBaseQuery("admins")
	query.Unions = []dbal.Union{{All: false, Query: unionQuery}}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "union")
	assert.Contains(t, result, "select * from (")
}

func TestCompileSelectUnionAll(t *testing.T) {
	g := newTestSQLite3WithQuoter()
	query := newBaseQuery("users")
	unionQuery := newBaseQuery("admins")
	query.Unions = []dbal.Union{{All: true, Query: unionQuery}}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "union all")
}

func TestCompileSelectWithColumns(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Columns = []interface{}{dbal.NewName("id"), dbal.NewName("name")}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "`id`")
	assert.Contains(t, result, "`name`")
}

func TestCompileSelectWithWhere(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Wheres = []dbal.Where{
		{Type: "basic", Column: "status", Operator: "=", Value: "active", Boolean: "and", Offset: 1},
	}
	query.Bindings["where"] = []interface{}{"active"}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "where `status` = ?")
}

func TestCompileSelectWithLimitAndOffset(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Limit = 10
	query.Offset = 5
	result := g.CompileSelect(query)
	assert.Contains(t, result, "limit")
	assert.Contains(t, result, "offset")
}

func TestCompileSelectDistinct(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Distinct = true
	result := g.CompileSelect(query)
	assert.Contains(t, result, "select")
}

// ---------------------------------------------------------------------------
// CompileWheres - IsJoinClause path
// ---------------------------------------------------------------------------

func TestCompileWheresJoinClause(t *testing.T) {
	g := newTestSQLite3()
	query := &dbal.Query{
		IsJoinClause: true,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	wheres := []dbal.Where{
		{Type: "basic", Column: "users.id", Operator: "=", Value: "orders.user_id", Boolean: "and", Offset: 1},
	}
	offset := 0
	result := g.CompileWheres(query, wheres, &offset)
	assert.Contains(t, result, "on ")
	assert.NotContains(t, result, "where ")
}

func TestCompileWheresEmptyReturnsEmpty(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	offset := 0
	result := g.CompileWheres(query, []dbal.Where{}, &offset)
	assert.Equal(t, "", result)
}

// ---------------------------------------------------------------------------
// WhereDate / WhereTime / WhereDay / WhereMonth / WhereYear
// ---------------------------------------------------------------------------

func TestWhereDate(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{Column: "created_at", Operator: "=", Value: "2024-01-01", Boolean: "and", Offset: 1}
	result := g.WhereDate(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%Y-%m-%d',`created_at`)")
	assert.Contains(t, result, "cast(? as text)")
	assert.Equal(t, 1, offset)
}

func TestWhereTime(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{Column: "created_at", Operator: "=", Value: "12:00:00", Boolean: "and", Offset: 1}
	result := g.WhereTime(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%H:%M:%S',`created_at`)")
	assert.Contains(t, result, "cast(? as text)")
}

func TestWhereDay(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{Column: "created_at", Operator: "=", Value: "15", Boolean: "and", Offset: 1}
	result := g.WhereDay(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%d',`created_at`)")
	assert.Contains(t, result, "cast(? as text)")
}

func TestWhereMonth(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{Column: "created_at", Operator: "=", Value: "06", Boolean: "and", Offset: 1}
	result := g.WhereMonth(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%m',`created_at`)")
	assert.Contains(t, result, "cast(? as text)")
}

func TestWhereYear(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{Column: "created_at", Operator: "=", Value: "2024", Boolean: "and", Offset: 1}
	result := g.WhereYear(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%Y',`created_at`)")
	assert.Contains(t, result, "cast(? as text)")
}

// ---------------------------------------------------------------------------
// WhereDateBased - expression value path
// ---------------------------------------------------------------------------

func TestWhereDateBasedWithExpression(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Column:   "created_at",
		Operator: "=",
		Value:    dbal.NewExpression("CURRENT_DATE"),
		Boolean:  "and",
		Offset:   1,
	}
	result := g.WhereDateBased("%Y-%m-%d", &dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%Y-%m-%d',`created_at`)")
	assert.Contains(t, result, "cast(CURRENT_DATE as text)")
	assert.Equal(t, 0, offset, "offset should not increment for expressions")
}

func TestWhereDateBasedWithNonExpression(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Column:   "updated_at",
		Operator: ">",
		Value:    "2024-06-01",
		Boolean:  "and",
		Offset:   1,
	}
	result := g.WhereDateBased("%Y-%m-%d", &dbal.Query{}, where, &offset)
	assert.Contains(t, result, "strftime('%Y-%m-%d',`updated_at`) > cast(? as text)")
	assert.Equal(t, 1, offset)
}

// ---------------------------------------------------------------------------
// WhereNested - IsJoinClause path
// ---------------------------------------------------------------------------

func TestWhereNestedJoinClause(t *testing.T) {
	g := newTestSQLite3()
	parentQuery := &dbal.Query{
		IsJoinClause: true,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	nestedQuery := &dbal.Query{
		Wheres: []dbal.Where{
			{Type: "basic", Column: "a", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
		},
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {1}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	where := dbal.Where{
		Type:    "nested",
		Boolean: "and",
		Query:   nestedQuery,
	}
	offset := 0
	result := g.WhereNested(parentQuery, where, &offset)
	assert.Contains(t, result, "(")
	assert.Contains(t, result, ")")
	assert.Contains(t, result, "`a` = ?")
}

func TestWhereNestedNormalQuery(t *testing.T) {
	g := newTestSQLite3()
	parentQuery := newBaseQuery("users")
	nestedQuery := &dbal.Query{
		Wheres: []dbal.Where{
			{Type: "basic", Column: "x", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
			{Type: "basic", Column: "y", Operator: ">", Value: 2, Boolean: "or", Offset: 1},
		},
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {1, 2}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	where := dbal.Where{Type: "nested", Boolean: "and", Query: nestedQuery}
	offset := 0
	result := g.WhereNested(parentQuery, where, &offset)
	assert.Contains(t, result, "(`x` = ?")
	assert.Contains(t, result, "or `y` > ?")
}

// ---------------------------------------------------------------------------
// CompileLock
// ---------------------------------------------------------------------------

func TestCompileLockReturnsEmpty(t *testing.T) {
	g := newTestSQLite3()
	result := g.CompileLock(&dbal.Query{}, "for update")
	assert.Equal(t, "", result)
}

// ---------------------------------------------------------------------------
// CompileInsertOrIgnore
// ---------------------------------------------------------------------------

func TestCompileInsertOrIgnore(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	columns := []interface{}{"name", "email"}
	values := [][]interface{}{{"John", "john@test.com"}}
	sql, bindings := g.CompileInsertOrIgnore(query, columns, values)
	assert.Contains(t, sql, "insert or ignore into")
	assert.Contains(t, sql, "`users`")
	assert.Equal(t, []interface{}{"John", "john@test.com"}, bindings)
}

func TestCompileInsertOrIgnoreMultipleRows(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	columns := []interface{}{"name"}
	values := [][]interface{}{{"Alice"}, {"Bob"}}
	sql, bindings := g.CompileInsertOrIgnore(query, columns, values)
	assert.Contains(t, sql, "insert or ignore into")
	assert.Equal(t, []interface{}{"Alice", "Bob"}, bindings)
}

// ---------------------------------------------------------------------------
// CompileUpsert
// ---------------------------------------------------------------------------

func TestCompileUpsertEmptyValues(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	sql, bindings := g.CompileUpsert(query, nil, [][]interface{}{}, []interface{}{"id"}, []interface{}{"name"})
	assert.Contains(t, sql, "insert into")
	assert.Contains(t, sql, "default values")
	assert.Empty(t, bindings)
}

func TestCompileUpsertWithSliceUpdateValues(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	columns := []interface{}{"id", "name", "email"}
	values := [][]interface{}{{1, "John", "john@test.com"}}
	uniqueBy := []interface{}{"id"}
	updateValues := []interface{}{"name", "email"}

	sql, bindings := g.CompileUpsert(query, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "do update set")
	assert.Contains(t, sql, "excluded")
	assert.Contains(t, sql, "`name`=excluded.`name`")
	assert.Contains(t, sql, "`email`=excluded.`email`")
	assert.Equal(t, []interface{}{1, "John", "john@test.com"}, bindings)
}

func TestCompileUpsertWithMapUpdateValues(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	columns := []interface{}{"id", "name"}
	values := [][]interface{}{{1, "John"}}
	uniqueBy := []interface{}{"id"}
	updateValues := map[string]interface{}{"name": "Jane"}

	sql, bindings := g.CompileUpsert(query, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "do update set")
	assert.Contains(t, sql, "`name`=?")
	assert.Contains(t, bindings, "Jane")
}

func TestCompileUpsertWithMapExpressionValue(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("counters")
	columns := []interface{}{"id", "count"}
	values := [][]interface{}{{1, 0}}
	uniqueBy := []interface{}{"id"}
	updateValues := map[string]interface{}{"count": dbal.NewExpression("count + 1")}

	sql, bindings := g.CompileUpsert(query, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "`count`=count + 1")
	assert.Equal(t, []interface{}{1, 0}, bindings)
}

// ---------------------------------------------------------------------------
// CompileTruncate
// ---------------------------------------------------------------------------

func TestCompileTruncate(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	sqls, bindings := g.CompileTruncate(query)
	assert.Len(t, sqls, 2)
	assert.Equal(t, "delete from sqlite_sequence where name = ?", sqls[0])
	assert.Contains(t, sqls[1], "delete from")
	assert.Contains(t, sqls[1], "`users`")
	assert.Len(t, bindings, 2)
	assert.Equal(t, []interface{}{"users"}, bindings[0])
	assert.Empty(t, bindings[1])
}

// ---------------------------------------------------------------------------
// getTableName
// ---------------------------------------------------------------------------

func TestGetTableNameWithName(t *testing.T) {
	g := newTestSQLite3()
	query := &dbal.Query{From: dbal.From{Name: dbal.NewName("orders")}}
	result := g.getTableName(query)
	assert.Equal(t, "orders", result)
}

func TestGetTableNameWithExpression(t *testing.T) {
	g := newTestSQLite3()
	query := &dbal.Query{From: dbal.From{Name: dbal.NewExpression("raw_table")}}
	result := g.getTableName(query)
	assert.Equal(t, "raw_table", result)
}

func TestGetTableNameWithString(t *testing.T) {
	g := newTestSQLite3()
	query := &dbal.Query{From: dbal.From{Name: "string_table"}}
	result := g.getTableName(query)
	assert.Equal(t, "string_table", result)
}

func TestGetTableNameWithDefaultType(t *testing.T) {
	g := newTestSQLite3()
	query := &dbal.Query{From: dbal.From{Name: 12345}}
	result := g.getTableName(query)
	assert.Equal(t, "12345", result)
}

// ---------------------------------------------------------------------------
// WrapUnion (quoter.go)
// ---------------------------------------------------------------------------

func TestWrapUnion(t *testing.T) {
	quoter := &Quoter{}
	result := quoter.WrapUnion("select * from users")
	assert.Equal(t, "select * from (select * from users)", result)
}

// ---------------------------------------------------------------------------
// CompileSelect additional edge cases
// ---------------------------------------------------------------------------

func TestCompileSelectWithGroupByAndHaving(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("orders")
	query.Groups = []interface{}{dbal.NewName("status")}
	query.Havings = []dbal.Having{
		{Type: "basic", Column: g.Raw("count(*)"), Operator: ">", Value: 5, Boolean: "and", Offset: 1},
	}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "group by")
	assert.Contains(t, result, "`status`")
}

func TestCompileSelectWithOrdering(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	query.Orders = []dbal.Order{{Column: "name", Direction: "asc"}}
	result := g.CompileSelect(query)
	assert.Contains(t, result, "order by")
}

func TestCompileSelectColumnsResetAfterCompile(t *testing.T) {
	g := newTestSQLite3()
	query := newBaseQuery("users")
	originalColumns := query.Columns
	g.CompileSelect(query)
	assert.Equal(t, originalColumns, query.Columns)
}

func TestGetOperatorsSQLite(t *testing.T) {
	g := newTestSQLite3()
	ops := g.GetOperators()
	assert.True(t, len(ops) > 5)
	assert.Contains(t, ops, "like")
	assert.Contains(t, ops, "ilike")
}
