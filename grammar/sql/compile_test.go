package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// ---------------------------------------------------------------------------
// Parameter / Parameterize nil handling (sql quoter)
// ---------------------------------------------------------------------------

func TestParameterPlainNil(t *testing.T) {
	g := newTestSQL()
	result := g.Parameter(nil, 1)
	assert.Equal(t, "NULL", result)
}

func TestParameterTypedNil(t *testing.T) {
	g := newTestSQL()
	var p *string
	result := g.Parameter(p, 1)
	assert.Equal(t, "NULL", result)
}

func TestParameterNonNil(t *testing.T) {
	g := newTestSQL()
	result := g.Parameter("hello", 1)
	assert.Equal(t, "?", result)
}

// ---------------------------------------------------------------------------
// CompileInsert nil handling
// ---------------------------------------------------------------------------

func TestCompileInsertWithNilValue(t *testing.T) {
	g := newTestSQL()
	q := &dbal.Query{From: dbal.From{Name: dbal.NewName("users")}}
	cols := []interface{}{"name", "deleted_at"}
	vals := [][]interface{}{{"alice", nil}}
	sql, bindings := g.CompileInsert(q, cols, vals)
	assert.Contains(t, sql, "NULL")
	assert.Equal(t, []interface{}{"alice"}, bindings)
}

func TestCompileInsertWithTypedNil(t *testing.T) {
	g := newTestSQL()
	q := &dbal.Query{From: dbal.From{Name: dbal.NewName("users")}}
	var p *string
	cols := []interface{}{"name", "deleted_at"}
	vals := [][]interface{}{{"bob", p}}
	sql, bindings := g.CompileInsert(q, cols, vals)
	assert.Contains(t, sql, "NULL")
	assert.Equal(t, []interface{}{"bob"}, bindings)
}

// ---------------------------------------------------------------------------
// CompileUpdateColumns nil handling
// ---------------------------------------------------------------------------

func TestCompileUpdateColumnsWithNil(t *testing.T) {
	g := newTestSQL()
	q := &dbal.Query{From: dbal.From{Name: dbal.NewName("users")}}
	offset := 0
	cols, bindings := g.CompileUpdateColumns(q, map[string]interface{}{"deleted_at": nil}, &offset)
	assert.Contains(t, cols, "NULL")
	assert.Empty(t, bindings)
	assert.Equal(t, 0, offset)
}

func TestCompileUpdateColumnsWithTypedNil(t *testing.T) {
	g := newTestSQL()
	q := &dbal.Query{From: dbal.From{Name: dbal.NewName("users")}}
	var p *string
	offset := 0
	cols, bindings := g.CompileUpdateColumns(q, map[string]interface{}{"deleted_at": p}, &offset)
	assert.Contains(t, cols, "NULL")
	assert.Empty(t, bindings)
	assert.Equal(t, 0, offset)
}

// ---------------------------------------------------------------------------
// WhereIn offset with nil
// ---------------------------------------------------------------------------

func TestWhereInOffsetWithNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: []interface{}{1, nil, 3},
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "in")
	assert.Equal(t, 2, offset, "nil should not count in offset")
}

func TestWhereInOffsetAllNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: []interface{}{nil, nil},
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Equal(t, 0, offset, "all nil -> offset stays 0")
}

func TestWhereInOffsetNoNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: []interface{}{1, 2, 3},
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Equal(t, 3, offset, "no nil -> offset = len(values)")
}

// ---------------------------------------------------------------------------
// WhereBetween offset with nil
// ---------------------------------------------------------------------------

func TestWhereBetweenOffsetMinNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{nil, 100},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "between")
	assert.Contains(t, result, "NULL")
	assert.Equal(t, 1, offset, "only max counts")
}

func TestWhereBetweenOffsetMaxNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{10, nil},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "between")
	assert.Contains(t, result, "NULL")
	assert.Equal(t, 1, offset, "only min counts")
}

func TestWhereBetweenOffsetBothNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{nil, nil},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Equal(t, 0, offset, "both nil -> offset unchanged")
}

func TestWhereBetweenOffsetNormal(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{10, 100},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Equal(t, 2, offset, "normal -> offset = 2")
}

// ---------------------------------------------------------------------------
// HavingBetween offset with nil
// ---------------------------------------------------------------------------

func TestHavingBetweenOffsetMinNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{nil, 1000},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.HavingBetween(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "between")
	assert.Contains(t, result, "NULL")
	assert.Equal(t, 1, offset)
}

func TestHavingBetweenOffsetBothNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{nil, nil},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	g.HavingBetween(&dbal.Query{}, having, &offset)
	assert.Equal(t, 0, offset)
}

func TestHavingBetweenOffsetNormal(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{100, 1000},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	g.HavingBetween(&dbal.Query{}, having, &offset)
	assert.Equal(t, 2, offset)
}

// ---------------------------------------------------------------------------
// HavingBasic offset with nil
// ---------------------------------------------------------------------------

func TestHavingBasicOffsetNil(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:     "basic",
		Column:   "total",
		Operator: "=",
		Value:    nil,
		Boolean:  "and",
		Offset:   1,
	}
	result := g.HavingBasic(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "NULL")
	assert.Equal(t, 0, offset, "nil value -> no offset increment")
}

func TestHavingBasicOffsetTypedNil(t *testing.T) {
	g := newTestSQL()
	var p *int
	offset := 0
	having := dbal.Having{
		Type:     "basic",
		Column:   "total",
		Operator: "=",
		Value:    p,
		Boolean:  "and",
		Offset:   1,
	}
	result := g.HavingBasic(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "NULL")
	assert.Equal(t, 0, offset)
}

func TestHavingBasicOffsetNormal(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:     "basic",
		Column:   "total",
		Operator: ">",
		Value:    100,
		Boolean:  "and",
		Offset:   1,
	}
	g.HavingBasic(&dbal.Query{}, having, &offset)
	assert.Equal(t, 1, offset)
}

// ---------------------------------------------------------------------------
// CompileHaving null / notnull types
// ---------------------------------------------------------------------------

func TestCompileHavingNull(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "null",
		Column:  "total",
		Boolean: "and",
	}
	result := g.CompileHaving(&dbal.Query{}, having, &offset)
	assert.Equal(t, "and `total` is null", result)
	assert.Equal(t, 0, offset)
}

func TestCompileHavingNotnull(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "notnull",
		Column:  "total",
		Boolean: "and",
	}
	result := g.CompileHaving(&dbal.Query{}, having, &offset)
	assert.Equal(t, "and `total` is not null", result)
	assert.Equal(t, 0, offset)
}

// ---------------------------------------------------------------------------
// utils.IsNil coverage for typed nil in parameter
// ---------------------------------------------------------------------------

func TestIsNilPlainNil(t *testing.T) {
	assert.True(t, utils.IsNil(nil))
}

func TestIsNilTypedNilPtr(t *testing.T) {
	var p *string
	assert.True(t, utils.IsNil(p))
}

func TestIsNilNonNilValue(t *testing.T) {
	s := "hello"
	assert.False(t, utils.IsNil(&s))
}

func TestIsNilNonPointer(t *testing.T) {
	assert.False(t, utils.IsNil(42))
}

// ---------------------------------------------------------------------------
// CompileHaving raw type
// ---------------------------------------------------------------------------

func TestCompileHavingRaw(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "raw",
		SQL:     "count(*) > 5",
		Boolean: "and",
	}
	result := g.CompileHaving(&dbal.Query{}, having, &offset)
	assert.Equal(t, "and count(*) > 5", result)
	assert.Equal(t, 0, offset)
}

// ---------------------------------------------------------------------------
// CompileHaving default branch (basic)
// ---------------------------------------------------------------------------

func TestCompileHavingBasicViaCompileHaving(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:     "basic",
		Column:   "total",
		Operator: ">",
		Value:    100,
		Boolean:  "and",
		Offset:   1,
	}
	result := g.CompileHaving(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, ">")
	assert.Equal(t, 1, offset)
}

// ---------------------------------------------------------------------------
// HavingBetween: Not=true branch
// ---------------------------------------------------------------------------

func TestHavingBetweenNotBetween(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{10, 1000},
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}
	result := g.HavingBetween(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "not between")
	assert.Equal(t, 2, offset)
}

// ---------------------------------------------------------------------------
// HavingBetween: Expression values
// ---------------------------------------------------------------------------

func TestHavingBetweenWithExpression(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{dbal.Raw("NOW()"), 1000},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.HavingBetween(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "NOW()")
	assert.Equal(t, 1, offset, "expression min skips offset, only max counts")
}

// ---------------------------------------------------------------------------
// WhereBetween: Not=true branch
// ---------------------------------------------------------------------------

func TestWhereBetweenNotBetween(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{10, 100},
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}
	result := g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "not between")
	assert.Equal(t, 2, offset)
}

// ---------------------------------------------------------------------------
// WhereBetween: Expression value
// ---------------------------------------------------------------------------

func TestWhereBetweenWithExpression(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{dbal.Raw("MIN(age)"), 100},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.WhereBetween(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "MIN(age)")
	assert.Equal(t, 1, offset, "expression min skips, only max counted")
}

// ---------------------------------------------------------------------------
// WhereIn: Not=true, nil ValuesIn, Expression subquery
// ---------------------------------------------------------------------------

func TestWhereInNotIn(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: []interface{}{1, 2, 3},
		Boolean:  "and",
		Not:      true,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "not in")
	assert.Equal(t, 3, offset)
}

func TestWhereInNilValues(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: nil,
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Equal(t, "false = true", result)
}

func TestWhereInNilValuesNot(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: nil,
		Boolean:  "and",
		Not:      true,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Equal(t, "true = true", result)
}

func TestWhereInExpressionSubquery(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: dbal.Raw("SELECT id FROM users"),
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "SELECT id FROM users")
}

func TestWhereInWithExpressionValues(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:     "in",
		Column:   "id",
		ValuesIn: []interface{}{1, dbal.Raw("NOW()"), 3},
		Boolean:  "and",
		Not:      false,
		Offset:   0,
	}
	result := g.WhereIn(&dbal.Query{}, where, &offset)
	assert.Contains(t, result, "NOW()")
	assert.Equal(t, 2, offset, "expression skipped in count")
}

// ---------------------------------------------------------------------------
// Parameter: Expression branch
// ---------------------------------------------------------------------------

func TestParameterExpression(t *testing.T) {
	g := newTestSQL()
	result := g.Parameter(dbal.Raw("NOW()"), 1)
	assert.Equal(t, "NOW()", result)
}

// ---------------------------------------------------------------------------
// CompileUpdateColumns: Expression value branch
// ---------------------------------------------------------------------------

func TestCompileUpdateColumnsWithExpression(t *testing.T) {
	g := newTestSQL()
	offset := 0
	query := &dbal.Query{}
	values := map[string]interface{}{
		"updated_at": dbal.Raw("NOW()"),
		"name":       "alice",
	}
	sql, bindings := g.CompileUpdateColumns(query, values, &offset)
	assert.Contains(t, sql, "NOW()")
	assert.Contains(t, sql, "?")
	assert.Equal(t, 1, len(bindings), "expression excluded from bindings")
}

// ---------------------------------------------------------------------------
// CompileHaving via CompileHaving: between branch
// ---------------------------------------------------------------------------

func TestCompileHavingBetweenBranch(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{10, 1000},
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result := g.CompileHaving(&dbal.Query{}, having, &offset)
	assert.Contains(t, result, "between")
	assert.Equal(t, 2, offset)
}

// ---------------------------------------------------------------------------
// HavingBetween: panic on invalid values length
// ---------------------------------------------------------------------------

func TestHavingBetweenPanicsOnBadLength(t *testing.T) {
	g := newTestSQL()
	offset := 0
	having := dbal.Having{
		Type:    "between",
		Column:  "total",
		Values:  []interface{}{10},
		Boolean: "and",
		Offset:  1,
	}
	assert.Panics(t, func() {
		g.HavingBetween(&dbal.Query{}, having, &offset)
	})
}

// ---------------------------------------------------------------------------
// WhereBetween: panic on invalid values length
// ---------------------------------------------------------------------------

func TestWhereBetweenPanicsOnBadLength(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "between",
		Column:  "age",
		Values:  []interface{}{10},
		Boolean: "and",
		Offset:  1,
	}
	assert.Panics(t, func() {
		g.WhereBetween(&dbal.Query{}, where, &offset)
	})
}

// ---------------------------------------------------------------------------
// CompileInsert: Expression, nil, empty values
// ---------------------------------------------------------------------------

func TestCompileInsertWithExpression(t *testing.T) {
	g := newTestSQL()
	query := newBaseQuery("users")
	columns := []interface{}{"name", "created_at"}
	values := [][]interface{}{{"alice", dbal.Raw("NOW()")}}
	sql, bindings := g.CompileInsert(query, columns, values)
	assert.Contains(t, sql, "NOW()")
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, "alice", bindings[0])
	_ = sql
}

func TestCompileInsertDefaultValues(t *testing.T) {
	g := newTestSQL()
	query := newBaseQuery("users")
	sql, bindings := g.CompileInsert(query, []interface{}{}, [][]interface{}{})
	assert.Contains(t, sql, "default values")
	assert.Nil(t, bindings)
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

func newTestSQL() SQL {
	return NewSQL(&Quoter{})
}

func TestWhereJsoncontainsMySQL(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsMySQLNot(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "not JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsMySQLOffset(t *testing.T) {
	g := newTestSQL()
	offset := 3
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"test"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 4, offset)
}
