package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	goSQL "github.com/yaoapp/xun/grammar/sql"
)

func newTestMySQL() MySQL {
	return MySQL{
		SQL: goSQL.NewSQL(&goSQL.Quoter{}),
	}
}

func newMySQLQuery(tableName string) *dbal.Query {
	return &dbal.Query{
		From:  dbal.From{Name: dbal.NewName(tableName)},
		Limit: -1,
		Bindings: map[string][]interface{}{
			"select": {}, "from": {}, "join": {},
			"where": {}, "groupBy": {}, "having": {}, "order": {},
		},
	}
}

func TestCompileUpsertMapWithNilMySQL(t *testing.T) {
	g := newTestMySQL()
	q := newMySQLQuery("users")
	columns := []interface{}{"name", "deleted_at"}
	values := [][]interface{}{{"alice", nil}}
	uniqueBy := []interface{}{"name"}
	updateValues := map[string]interface{}{"deleted_at": nil}

	sql, bindings := g.CompileUpsert(q, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on duplicate key update")
	assert.Contains(t, sql, "NULL")
	assert.Len(t, bindings, 1, "nil excluded from bindings")
}

func TestCompileUpsertMapWithTypedNilMySQL(t *testing.T) {
	g := newTestMySQL()
	q := newMySQLQuery("users")
	var p *string
	columns := []interface{}{"name", "deleted_at"}
	values := [][]interface{}{{"bob", p}}
	uniqueBy := []interface{}{"name"}
	updateValues := map[string]interface{}{"deleted_at": p}

	sql, bindings := g.CompileUpsert(q, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on duplicate key update")
	assert.Contains(t, sql, "NULL")
	assert.Len(t, bindings, 1, "typed nil excluded from bindings")
}

func TestCompileUpsertMapNormalValueMySQL(t *testing.T) {
	g := newTestMySQL()
	q := newMySQLQuery("users")
	columns := []interface{}{"name", "email"}
	values := [][]interface{}{{"alice", "alice@test.com"}}
	uniqueBy := []interface{}{"email"}
	updateValues := map[string]interface{}{"name": "alice_updated"}

	sql, bindings := g.CompileUpsert(q, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on duplicate key update")
	assert.Len(t, bindings, 3)
}

func TestCompileUpsertSliceMySQL(t *testing.T) {
	g := newTestMySQL()
	q := newMySQLQuery("users")
	columns := []interface{}{"name", "email"}
	values := [][]interface{}{{"alice", "alice@test.com"}}
	uniqueBy := []interface{}{"email"}
	updateValues := []string{"name"}

	sql, bindings := g.CompileUpsert(q, columns, values, uniqueBy, updateValues)
	assert.Contains(t, sql, "on duplicate key update")
	assert.Contains(t, sql, "values(`name`)")
	assert.Len(t, bindings, 2)
}

func TestCompileUpsertEmptyMySQL(t *testing.T) {
	g := newTestMySQL()
	q := newMySQLQuery("users")
	sql, bindings := g.CompileUpsert(q, nil, [][]interface{}{}, nil, nil)
	assert.Contains(t, sql, "default values")
	assert.Empty(t, bindings)
}
