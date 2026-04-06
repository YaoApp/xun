package postgres

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	goSQL "github.com/yaoapp/xun/grammar/sql"
)

func newTestPostgres() Postgres {
	return Postgres{
		SQL: goSQL.NewSQL(&Quoter{}),
	}
}

func TestWhereJsoncontainsPG(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `"tags"::jsonb @> $1`, result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsPGNot(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `not "tags"::jsonb @> $1`, result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsPGOffset(t *testing.T) {
	pg := newTestPostgres()
	offset := 2
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"test"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `"tags"::jsonb @> $3`, result)
	assert.Equal(t, 3, offset)
}

func TestWhereJsoncontainsPGMixed(t *testing.T) {
	pg := newTestPostgres()
	offset := 0

	where1 := dbal.Where{
		Type:     "basic",
		Column:   "name",
		Operator: "=",
		Value:    "test",
		Boolean:  "and",
		Not:      false,
		Offset:   1,
	}
	result1 := pg.WhereBasic(&dbal.Query{}, where1, &offset)
	assert.Equal(t, `"name" = $1`, result1)
	assert.Equal(t, 1, offset)

	where2 := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result2 := pg.WhereJsoncontains(&dbal.Query{}, where2, &offset)
	assert.Equal(t, `"tags"::jsonb @> $2`, result2)
	assert.Equal(t, 2, offset)
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

func TestCompileDeleteSimplePG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("users", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := pg.CompileDelete(query)
	assert.Equal(t, `delete from "users" where "id" = $1`, sql)
	assert.Equal(t, []interface{}{1}, bindings)
}

func TestCompileDeleteWithJsonContainsPG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("assistants", []dbal.Where{
		{Type: "jsoncontains", Column: "tags", Value: `"test"`, Boolean: "and", Offset: 1},
	}, []interface{}{`"test"`})

	sql, bindings := pg.CompileDelete(query)
	assert.Equal(t, `delete from "assistants" where "tags"::jsonb @> $1`, sql)
	assert.Equal(t, []interface{}{`"test"`}, bindings)
}

func TestCompileDeleteNoWherePG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("users", nil, nil)
	query.Bindings["where"] = []interface{}{}

	sql, bindings := pg.CompileDelete(query)
	assert.Equal(t, `delete from "users" `, sql)
	assert.Empty(t, bindings)
}

func TestCompileUpdateSimplePG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("users", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := pg.CompileUpdate(query, map[string]interface{}{"name": "new_name"})
	assert.Contains(t, sql, `update "users" set`)
	assert.Contains(t, sql, `where "id" =`)
	assert.Contains(t, bindings, 1)
	assert.Contains(t, bindings, "new_name")
}

func TestCompileUpdateWithJsonContainsPG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("assistants", []dbal.Where{
		{Type: "jsoncontains", Column: "tags", Value: `"old"`, Boolean: "and", Offset: 1},
	}, []interface{}{`"old"`})

	sql, bindings := pg.CompileUpdate(query, map[string]interface{}{"name": "updated"})
	assert.Contains(t, sql, `update "assistants" set`)
	assert.Contains(t, sql, `"tags"::jsonb @>`)
	assert.Contains(t, bindings, `"old"`)
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

func TestCompileDeleteWithLimitPG(t *testing.T) {
	pg := newTestPostgres()
	query := newQueryWithLimit("users", 10, []dbal.Where{
		{Type: "basic", Column: "status", Operator: "=", Value: "inactive", Boolean: "and", Offset: 1},
	}, []interface{}{"inactive"})

	sql, bindings := pg.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, `"ctid"`)
	assert.Contains(t, sql, "in (")
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteWithLimitAliasPG(t *testing.T) {
	pg := newTestPostgres()
	query := newQueryWithAlias("users", "u", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := pg.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, `"ctid"`)
	assert.Contains(t, sql, `"u"."ctid"`)
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteWithLimitNoAliasPG(t *testing.T) {
	pg := newTestPostgres()
	query := newQueryWithLimit("users", 5, nil, nil)
	query.Bindings["where"] = []interface{}{}

	sql, _ := pg.CompileDelete(query)
	assert.Contains(t, sql, "delete from")
	assert.Contains(t, sql, `"ctid"`)
	assert.Contains(t, sql, "ctid")
}

func TestCompileUpdateWithLimitPG(t *testing.T) {
	pg := newTestPostgres()
	query := newQueryWithLimit("users", 10, []dbal.Where{
		{Type: "basic", Column: "status", Operator: "=", Value: "old", Boolean: "and", Offset: 1},
	}, []interface{}{"old"})

	sql, bindings := pg.CompileUpdate(query, map[string]interface{}{"status": "new"})
	assert.Contains(t, sql, "update")
	assert.Contains(t, sql, `"ctid"`)
	assert.Contains(t, sql, "in (")
	assert.NotEmpty(t, bindings)
}

func TestCompileUpdateWithLimitAliasPG(t *testing.T) {
	pg := newTestPostgres()
	query := newQueryWithAlias("users", "u", []dbal.Where{
		{Type: "basic", Column: "id", Operator: "=", Value: 1, Boolean: "and", Offset: 1},
	}, []interface{}{1})

	sql, bindings := pg.CompileUpdate(query, map[string]interface{}{"name": "updated"})
	assert.Contains(t, sql, "update")
	assert.Contains(t, sql, `"ctid"`)
	assert.Contains(t, sql, `"u"."ctid"`)
	assert.NotEmpty(t, bindings)
}

func TestCompileDeleteNestedJsonContainsPG(t *testing.T) {
	pg := newTestPostgres()
	query := newDeleteQuery("assistants", []dbal.Where{
		{
			Type:    "nested",
			Boolean: "and",
			Query: &dbal.Query{
				Wheres: []dbal.Where{
					{Type: "jsoncontains", Column: "tags", Value: `"a"`, Boolean: "and", Offset: 1},
					{Type: "jsoncontains", Column: "tags", Value: `"b"`, Boolean: "or", Offset: 1},
				},
				Bindings: map[string][]interface{}{
					"select": {}, "from": {}, "join": {},
					"where":   {`"a"`, `"b"`},
					"groupBy": {}, "having": {}, "order": {},
				},
			},
		},
	}, []interface{}{`"a"`, `"b"`})

	sql, bindings := pg.CompileDelete(query)
	assert.Contains(t, sql, `"tags"::jsonb @>`)
	assert.Contains(t, sql, "delete from")
	assert.Len(t, bindings, 2)
}

// --- helpers ---

func newFullQuery() *dbal.Query {
	q := dbal.NewQuery()
	q.From = dbal.From{Name: dbal.NewName("users")}
	return q
}

// --- CompileSelect / CompileSelectOffset ---

func TestCompileSelectBasic(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "select")
	assert.Contains(t, result, "*")
}

func TestCompileSelectWithSQLStmt(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.SQL = "select * from users where id = $1"
	result := pg.CompileSelect(q)
	assert.Equal(t, `select * from users where id = $1`, result)
}

func TestCompileSelectWithSQLStmtAppendsLimitOffset(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.SQL = "select * from users where active = true"
	q.Limit = 10
	q.Offset = 5
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "limit 10")
	assert.Contains(t, result, "offset 5")
}

func TestCompileSelectWithSQLStmtAlreadyHasLimit(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.SQL = "select * from users limit 5"
	result := pg.CompileSelect(q)
	assert.Equal(t, "select * from users limit 5", result)
}

func TestCompileSelectWithSQLStmtAlreadyHasOffset(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.SQL = "select * from users offset 10"
	result := pg.CompileSelect(q)
	assert.Equal(t, "select * from users offset 10", result)
}

func TestCompileSelectUnionAggregate(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Aggregate = dbal.Aggregate{Func: "count", Columns: []interface{}{pg.Raw("*")}}
	q.Unions = []dbal.Union{
		{All: false, Query: newFullQuery()},
	}
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "count")
}

func TestCompileSelectNormalPath(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Columns = []interface{}{dbal.NewName("id"), dbal.NewName("name")}
	q.Limit = 10
	q.Offset = 0
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "select")
	assert.Contains(t, result, `"id"`)
	assert.Contains(t, result, `"name"`)
	assert.Contains(t, result, "limit 10")
	assert.Contains(t, result, "offset 0")
}

func TestCompileSelectEmptyColumnsDefaultsStar(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Columns = []interface{}{}
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "*")
}

func TestCompileSelectWithUnions(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q2 := newFullQuery()
	q2.From = dbal.From{Name: dbal.NewName("admins")}
	q.Unions = []dbal.Union{
		{All: false, Query: q2},
	}
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "union")
}

func TestCompileSelectWithLock(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Lock = "share"
	result := pg.CompileSelect(q)
	assert.Contains(t, result, "for share")
}

// --- CompileColumns ---

func TestCompileColumnsWithAggregate(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Aggregate = dbal.Aggregate{Func: "count", Columns: []interface{}{pg.Raw("*")}}
	offset := 0
	result := pg.CompileColumns(q, q.Columns, &offset)
	assert.Equal(t, "", result)
}

func TestCompileColumnsDistinct(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.Distinct = true
	q.Columns = []interface{}{dbal.NewName("id")}
	offset := 0
	result := pg.CompileColumns(q, q.Columns, &offset)
	assert.Contains(t, result, "select distinct")
	assert.NotContains(t, result, "distinct on")
}

func TestCompileColumnsDistinctOn(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.DistinctColumns = []interface{}{dbal.NewName("email")}
	q.Columns = []interface{}{dbal.NewName("id"), dbal.NewName("email")}
	offset := 0
	result := pg.CompileColumns(q, q.Columns, &offset)
	assert.Contains(t, result, "select distinct on")
	assert.Contains(t, result, `"email"`)
}

func TestCompileColumnsWithSelectOffset(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{
		dbal.Select{SQL: "(select count(*) from orders)", Alias: "order_count", Offset: 2},
	}
	offset := 0
	pg.CompileColumns(q, cols, &offset)
	assert.Equal(t, 2, offset)
}

// --- CompileWheres with IsJoinClause ---

func TestCompileWheresJoinClause(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	q.IsJoinClause = true
	wheres := []dbal.Where{
		{Type: "basic", Column: "u.id", Operator: "=", Value: "o.user_id", Boolean: "and", Offset: 1},
	}
	offset := 0
	result := pg.CompileWheres(q, wheres, &offset)
	assert.True(t, strings.HasPrefix(result, "on "))
}

func TestCompileWheresEmpty(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	offset := 0
	result := pg.CompileWheres(q, []dbal.Where{}, &offset)
	assert.Equal(t, "", result)
}

// --- WhereDate ---

func TestWhereDateNonExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "date", Column: "created_at", Operator: "=", Value: "2024-01-01", Boolean: "and", Offset: 1,
	}
	result := pg.WhereDate(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, `"created_at"::date`)
	assert.Contains(t, result, "=$1")
	assert.Equal(t, 1, offset)
}

func TestWhereDateExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "date", Column: "created_at", Operator: "=",
		Value: dbal.NewExpression("CURRENT_DATE"), Boolean: "and", Offset: 1,
	}
	result := pg.WhereDate(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, `"created_at"::date`)
	assert.Contains(t, result, "CURRENT_DATE")
	assert.Equal(t, 0, offset)
}

// --- WhereTime ---

func TestWhereTimeNonExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "time", Column: "start_time", Operator: ">=", Value: "08:00:00", Boolean: "and", Offset: 1,
	}
	result := pg.WhereTime(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, `"start_time"::time`)
	assert.Contains(t, result, ">=$1")
}

func TestWhereTimeExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "time", Column: "start_time", Operator: "=",
		Value: dbal.NewExpression("CURRENT_TIME"), Boolean: "and", Offset: 1,
	}
	result := pg.WhereTime(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, `"start_time"::time`)
	assert.Contains(t, result, "CURRENT_TIME")
	assert.Equal(t, 0, offset)
}

// --- WhereDay / WhereMonth / WhereYear / WhereDateBased ---

func TestWhereDayNonExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "day", Column: "created_at", Operator: "=", Value: 15, Boolean: "and", Offset: 1,
	}
	result := pg.WhereDay(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(day from")
	assert.Contains(t, result, "$1")
}

func TestWhereDayExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "day", Column: "created_at", Operator: "=",
		Value: dbal.NewExpression("1"), Boolean: "and", Offset: 1,
	}
	result := pg.WhereDay(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(day from")
	assert.Contains(t, result, "1")
	assert.Equal(t, 0, offset)
}

func TestWhereMonthNonExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "month", Column: "created_at", Operator: "=", Value: 6, Boolean: "and", Offset: 1,
	}
	result := pg.WhereMonth(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(month from")
	assert.Contains(t, result, "$1")
}

func TestWhereMonthExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "month", Column: "created_at", Operator: "=",
		Value: dbal.NewExpression("6"), Boolean: "and", Offset: 1,
	}
	result := pg.WhereMonth(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(month from")
	assert.Contains(t, result, "6")
}

func TestWhereYearNonExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "year", Column: "created_at", Operator: "=", Value: 2024, Boolean: "and", Offset: 1,
	}
	result := pg.WhereYear(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(year from")
	assert.Contains(t, result, "$1")
}

func TestWhereYearExpression(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type: "year", Column: "created_at", Operator: "=",
		Value: dbal.NewExpression("2024"), Boolean: "and", Offset: 1,
	}
	result := pg.WhereYear(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "extract(year from")
	assert.Contains(t, result, "2024")
}

// --- WhereNested with IsJoinClause ---

func TestWhereNestedJoinClause(t *testing.T) {
	pg := newTestPostgres()
	inner := &dbal.Query{
		IsJoinClause: true,
		Wheres: []dbal.Where{
			{Type: "basic", Column: "a.id", Operator: "=", Value: "b.id", Boolean: "and", Offset: 1},
		},
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {"b.id"}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	outer := &dbal.Query{IsJoinClause: true}
	where := dbal.Where{Type: "nested", Boolean: "and", Query: inner}
	offset := 0
	result := pg.WhereNested(outer, where, &offset)
	assert.True(t, strings.HasPrefix(result, "("))
	assert.True(t, strings.HasSuffix(result, ")"))
}

func TestWhereNestedNormalQuery(t *testing.T) {
	pg := newTestPostgres()
	inner := &dbal.Query{
		Wheres: []dbal.Where{
			{Type: "basic", Column: "status", Operator: "=", Value: "active", Boolean: "and", Offset: 1},
			{Type: "basic", Column: "role", Operator: "=", Value: "admin", Boolean: "or", Offset: 1},
		},
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {"active", "admin"}, "groupBy": {}, "having": {}, "order": {},
		},
	}
	where := dbal.Where{Type: "nested", Boolean: "and", Query: inner}
	offset := 0
	result := pg.WhereNested(&dbal.Query{}, where, &offset)
	assert.True(t, strings.HasPrefix(result, "("))
	assert.True(t, strings.HasSuffix(result, ")"))
	assert.Contains(t, result, `"status"`)
}

// --- CompileLock ---

func TestCompileLockShare(t *testing.T) {
	pg := newTestPostgres()
	result := pg.CompileLock(&dbal.Query{}, "share")
	assert.Equal(t, "for share", result)
}

func TestCompileLockUpdate(t *testing.T) {
	pg := newTestPostgres()
	result := pg.CompileLock(&dbal.Query{}, "update")
	assert.Equal(t, "for update", result)
}

func TestCompileLockInvalid(t *testing.T) {
	pg := newTestPostgres()
	result := pg.CompileLock(&dbal.Query{}, "invalid")
	assert.Equal(t, "", result)
}

func TestCompileLockNonString(t *testing.T) {
	pg := newTestPostgres()
	result := pg.CompileLock(&dbal.Query{}, 123)
	assert.Equal(t, "", result)
}

func TestCompileLockNil(t *testing.T) {
	pg := newTestPostgres()
	result := pg.CompileLock(&dbal.Query{}, nil)
	assert.Equal(t, "", result)
}

// --- CompileInsertOrIgnore ---

func TestCompileInsertOrIgnore(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name", "email"}
	vals := [][]interface{}{{"alice", "alice@example.com"}}
	sql, bindings := pg.CompileInsertOrIgnore(q, cols, vals)
	assert.Contains(t, sql, "insert into")
	assert.Contains(t, sql, "on conflict do nothing")
	assert.Equal(t, []interface{}{"alice", "alice@example.com"}, bindings)
}

func TestCompileInsertOrIgnoreEmpty(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	sql, bindings := pg.CompileInsertOrIgnore(q, []interface{}{}, [][]interface{}{})
	assert.Contains(t, sql, "default values")
	assert.Contains(t, sql, "on conflict do nothing")
	assert.Nil(t, bindings)
}

// --- CompileInsertGetID ---

func TestCompileInsertGetID(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name"}
	vals := [][]interface{}{{"bob"}}
	sql, bindings := pg.CompileInsertGetID(q, cols, vals, "id")
	assert.Contains(t, sql, "insert into")
	assert.Contains(t, sql, `returning "id"`)
	assert.Equal(t, []interface{}{"bob"}, bindings)
}

func TestCompileInsertGetIDCustomSequence(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name"}
	vals := [][]interface{}{{"charlie"}}
	sql, _ := pg.CompileInsertGetID(q, cols, vals, "user_id")
	assert.Contains(t, sql, `returning "user_id"`)
}

// --- CompileUpsert ---

func TestCompileUpsertEmptyValues(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	sql, bindings := pg.CompileUpsert(q, nil, [][]interface{}{}, nil, nil)
	assert.Contains(t, sql, "default values")
	assert.Empty(t, bindings)
}

func TestCompileUpsertSliceUpdateValues(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name", "email"}
	vals := [][]interface{}{{"alice", "alice@example.com"}}
	uniqueBy := []interface{}{"email"}
	updateVals := []string{"name"}
	sql, bindings := pg.CompileUpsert(q, cols, vals, uniqueBy, updateVals)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "do update set")
	assert.Contains(t, sql, "excluded")
	assert.Len(t, bindings, 2)
}

func TestCompileUpsertMapUpdateValues(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name", "email"}
	vals := [][]interface{}{{"alice", "alice@example.com"}}
	uniqueBy := []interface{}{"email"}
	updateVals := map[string]interface{}{"name": "alice_updated"}
	sql, bindings := pg.CompileUpsert(q, cols, vals, uniqueBy, updateVals)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "do update set")
	assert.Contains(t, sql, "$")
	assert.Len(t, bindings, 3)
}

func TestCompileUpsertMapWithExpression(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name", "counter"}
	vals := [][]interface{}{{"alice", 1}}
	uniqueBy := []interface{}{"name"}
	updateVals := map[string]interface{}{"counter": dbal.NewExpression("counter + 1")}
	sql, bindings := pg.CompileUpsert(q, cols, vals, uniqueBy, updateVals)
	assert.Contains(t, sql, "counter + 1")
	assert.Len(t, bindings, 2)
}

func TestCompileUpsertMapWithNil(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	cols := []interface{}{"name", "deleted_at"}
	vals := [][]interface{}{{"alice", nil}}
	uniqueBy := []interface{}{"name"}
	updateVals := map[string]interface{}{"deleted_at": nil}
	sql, bindings := pg.CompileUpsert(q, cols, vals, uniqueBy, updateVals)
	assert.Contains(t, sql, "on conflict")
	assert.Contains(t, sql, "NULL")
	assert.Len(t, bindings, 1)
}

// --- CompileTruncate ---

func TestCompileTruncate(t *testing.T) {
	pg := newTestPostgres()
	q := newFullQuery()
	sqls, bindings := pg.CompileTruncate(q)
	assert.Len(t, sqls, 1)
	assert.Contains(t, sqls[0], "truncate table")
	assert.Contains(t, sqls[0], `"users"`)
	assert.Contains(t, sqls[0], "restart identity cascade")
	assert.Len(t, bindings, 1)
	assert.Empty(t, bindings[0])
}

// --- Quoter.VAL ---

func TestVALString(t *testing.T) {
	q := &Quoter{}
	result := q.VAL("hello")
	assert.Equal(t, "'hello'", result)
}

func TestVALStringPtr(t *testing.T) {
	q := &Quoter{}
	s := "world"
	result := q.VAL(&s)
	assert.Equal(t, "'world'", result)
}

func TestVALInt(t *testing.T) {
	q := &Quoter{}
	result := q.VAL(42)
	assert.Equal(t, "'42'", result)
}

func TestVALDefault(t *testing.T) {
	q := &Quoter{}
	result := q.VAL(true)
	assert.Equal(t, "'true'", result)
}

func TestVALEscapeQuotes(t *testing.T) {
	q := &Quoter{}
	result := q.VAL("it's")
	assert.Equal(t, `'it\'s'`, result)
}

func TestVALStripNewlines(t *testing.T) {
	q := &Quoter{}
	result := q.VAL("line1\nline2\r")
	assert.Equal(t, "'line1line2'", result)
}

// --- Quoter.Wrap ---

func TestWrapExpression(t *testing.T) {
	q := &Quoter{}
	expr := dbal.NewExpression("NOW()")
	result := q.Wrap(expr)
	assert.Equal(t, "NOW()", result)
}

func TestWrapNameWithAs(t *testing.T) {
	q := &Quoter{}
	name := dbal.NewName("users as u")
	result := q.Wrap(name)
	assert.Equal(t, `"users" as u`, result)
}

func TestWrapNameWithoutAs(t *testing.T) {
	q := &Quoter{}
	name := dbal.NewName("users")
	result := q.Wrap(name)
	assert.Equal(t, `"users"`, result)
}

func TestWrapSelectWithAlias(t *testing.T) {
	q := &Quoter{}
	sel := dbal.Select{SQL: "(select count(*) from orders)", Alias: "cnt"}
	result := q.Wrap(sel)
	assert.Equal(t, `(select count(*) from orders) as "cnt"`, result)
}

func TestWrapSelectWithoutAlias(t *testing.T) {
	q := &Quoter{}
	sel := dbal.Select{SQL: "(select 1)"}
	result := q.Wrap(sel)
	assert.Equal(t, "(select 1) ", result)
}

func TestWrapString(t *testing.T) {
	q := &Quoter{}
	result := q.Wrap("username")
	assert.Equal(t, `"username"`, result)
}

func TestWrapDefault(t *testing.T) {
	q := &Quoter{}
	result := q.Wrap(12345)
	assert.Equal(t, "12345", result)
}

// --- Quoter.WrapAliasedValue ---

func TestWrapAliasedValueStar(t *testing.T) {
	q := &Quoter{}
	result := q.WrapAliasedValue("*")
	assert.Equal(t, "*", result)
}

func TestWrapAliasedValueDotted(t *testing.T) {
	q := &Quoter{}
	result := q.WrapAliasedValue("users.name")
	assert.Equal(t, `"users"."name"`, result)
}

func TestWrapAliasedValueWithAlias(t *testing.T) {
	q := &Quoter{}
	result := q.WrapAliasedValue("username as uname")
	assert.Contains(t, result, `as`)
}

func TestWrapAliasedValuePlain(t *testing.T) {
	q := &Quoter{}
	result := q.WrapAliasedValue("email")
	assert.Equal(t, `"email"`, result)
}

// --- Quoter.WrapTable ---

func TestWrapTableExpression(t *testing.T) {
	q := &Quoter{}
	expr := dbal.NewExpression("raw_table")
	result := q.WrapTable(expr)
	assert.Equal(t, "raw_table", result)
}

func TestWrapTableNameWithAs(t *testing.T) {
	q := &Quoter{}
	name := dbal.NewName("users as u")
	result := q.WrapTable(name)
	assert.Equal(t, `"users" as "u"`, result)
}

func TestWrapTableNameWithoutAs(t *testing.T) {
	q := &Quoter{}
	name := dbal.NewName("users")
	result := q.WrapTable(name)
	assert.Equal(t, `"users"`, result)
}

func TestWrapTableFrom(t *testing.T) {
	q := &Quoter{}
	from := dbal.From{Name: dbal.NewName("products")}
	result := q.WrapTable(from)
	assert.Equal(t, `"products"`, result)
}

func TestWrapTableString(t *testing.T) {
	q := &Quoter{}
	result := q.WrapTable("orders")
	assert.Equal(t, `"orders"`, result)
}

func TestWrapTableDefault(t *testing.T) {
	q := &Quoter{}
	result := q.WrapTable(999)
	assert.Equal(t, "999", result)
}

// --- Quoter.Parameter ---

func TestParameterExpression(t *testing.T) {
	q := &Quoter{}
	expr := dbal.NewExpression("DEFAULT")
	result := q.Parameter(expr, 1)
	assert.Equal(t, "DEFAULT", result)
}

func TestParameterNil(t *testing.T) {
	q := &Quoter{}
	result := q.Parameter(nil, 1)
	assert.Equal(t, "NULL", result)
}

func TestParameterNormal(t *testing.T) {
	q := &Quoter{}
	result := q.Parameter("hello", 3)
	assert.Equal(t, "$3", result)
}

// --- Quoter.Parameterize ---

func TestParameterizeMixed(t *testing.T) {
	q := &Quoter{}
	values := []interface{}{"a", nil, dbal.NewExpression("NOW()"), "b"}
	result := q.Parameterize(values, 0)
	assert.Contains(t, result, "$1")
	assert.Contains(t, result, "NULL")
	assert.Contains(t, result, "NOW()")
	assert.Contains(t, result, "$2")
}

func TestParameterizeAllNormal(t *testing.T) {
	q := &Quoter{}
	values := []interface{}{"x", "y", "z"}
	result := q.Parameterize(values, 0)
	assert.Equal(t, "$1,$2,$3", result)
}

func TestParameterizeWithOffset(t *testing.T) {
	q := &Quoter{}
	values := []interface{}{"x", "y"}
	result := q.Parameterize(values, 5)
	assert.Equal(t, "$6,$7", result)
}

// --- Quoter.Columnize ---

func TestColumnize(t *testing.T) {
	q := &Quoter{}
	cols := []interface{}{"id", "name", "email"}
	result := q.Columnize(cols)
	assert.Contains(t, result, `"id"`)
	assert.Contains(t, result, `"name"`)
	assert.Contains(t, result, `"email"`)
}

func TestColumnizeWithExpression(t *testing.T) {
	q := &Quoter{}
	cols := []interface{}{"id", dbal.NewExpression("count(*)")}
	result := q.Columnize(cols)
	assert.Contains(t, result, `"id"`)
	assert.Contains(t, result, "count(*)")
}

// --- Quoter.ID ---

func TestIDBasic(t *testing.T) {
	q := Quoter{}
	result := q.ID("users")
	assert.Equal(t, `"users"`, result)
}

func TestIDStripsQuotes(t *testing.T) {
	q := Quoter{}
	result := q.ID(`"users"`)
	assert.Equal(t, `"users"`, result)
}

func TestIDStripsNewlines(t *testing.T) {
	q := Quoter{}
	result := q.ID("us\ners\r")
	assert.Equal(t, `"users"`, result)
}

func TestGetOperatorsPG(t *testing.T) {
	pg := newTestPostgres()
	ops := pg.GetOperators()
	assert.True(t, len(ops) > 10)
	assert.Contains(t, ops, "@>")
	assert.Contains(t, ops, "ilike")
}
