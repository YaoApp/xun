package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestUnionUnionAll(t *testing.T) {
	NewTableFoUnionTest()
	qb := getTestBuilder()
	qb.Table("table_test_union_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		UnionAll(func(qb Query) {
			qb.Table("table_test_union_t2").
				Where("email", "like", "%yaojs.org").
				Select("id", "email", "name")
		})

	// MySQL: (select `id`, `email`, `name` from `table_test_union_t1` where `email` like ? ) union all (select `id`, `email`, `name` from `table_test_union_t2` where `email` like ?)
	// Postgres: (select "id", "email", "name" from "table_test_union_t1" where "email" like $1 ) union all (select "id", "email", "name" from "table_test_union_t2" where "email" like $2)
	// SQLite3: select * from (select `id`, `email`, `name` from `table_test_union_t1` where `email` like ? ) union all select * from (select `id`, `email`, `name` from `table_test_union_t2` where `email` like ?)

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `(select "id", "email", "name" from "table_test_union_t1" where "email" like $1 ) union all (select "id", "email", "name" from "table_test_union_t2" where "email" like $2)`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select * from (select `id`, `email`, `name` from `table_test_union_t1` where `email` like ? ) union all select * from (select `id`, `email`, `name` from `table_test_union_t2` where `email` like ?)", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "(select `id`, `email`, `name` from `table_test_union_t1` where `email` like ? ) union all (select `id`, `email`, `name` from `table_test_union_t2` where `email` like ?)", sql, "the query sql not equal")
	}

	bindings := qb.GetBindings()
	assert.Equal(t, 2, len(bindings), "the bindings should have 2 items")
	if len(bindings) == 1 {
		assert.Equal(t, "%@yao.run", bindings[0].(string), "the 1st binding should be %@yao.run")
		assert.Equal(t, "%@yaojs.org", bindings[1].(string), "the 2nd binding should be %@yaojs.org")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 7, len(rows), "the return value should has 7 rows")
	if len(rows) == 7 {
		assert.Equal(t, "john@yao.run", rows[0]["email"].(string), "the email of first row should be john@yao.run")
		assert.Equal(t, "nio@yaojs.org", rows[4]["email"].(string), "the email of 4th row should be nio@yaojs.org")
	}
}

// @todo: test union

// @todo: test unionOrders

// @todo: test unionLimit

// @todo: test unionOffset

// clean the test data
func TestUnionClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_union")
}

func NewTableFoUnionTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_union_t1")
	builder.MustCreateTable("table_test_union_t1", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email")
		table.String("name")
		table.Integer("vote")
		table.Float("score", 5, 2)
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})

	qb := getTestBuilder()
	qb.Table("table_test_union_t1").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})

	builder.DropTableIfExists("table_test_union_t2")
	builder.MustCreateTable("table_test_union_t2", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email")
		table.String("name")
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_union_t2").Insert([]xun.R{
		{"email": "nio@yaojs.org", "name": "Nio", "status": "WAITING", "created_at": "2021-03-26 00:15:16"},
		{"email": "Tom@yaojs.org", "name": "Tom", "status": "PENDING", "created_at": "2021-03-26 08:19:15"},
		{"email": "Han@yaojs.org", "name": "Han", "status": "DONE", "created_at": "2021-03-26 10:24:23"},
	})
}
