package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestFromFrom(t *testing.T) {
	NewTableFoFromTest()
	qb := getTestBuilder()
	qb.From("table_test_from as t").
		Where("email", "like", "%@yao.run").
		Select("id", "t.email as wid", "t.cate as category").
		OrderBy("id")

	// select `id`, `t`.`email` as `wid`, `t`.`cate` as `category` from `table_test_from` as `t` where `email` like ? order by `id` asc
	// select "id", "t"."email" as "wid", "t"."cate" as "category" from "table_test_from" as "t" where "email" like $1 order by "id" asc
	checkingFromFrom(t, qb)
}

func TestFromFromRaw(t *testing.T) {
	NewTableFoFromTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.FromRaw(`"table_test_from" as "t"`).
			Where("email", "like", "%@yao.run").
			Select("id", "t.email as wid", "t.cate as category").
			OrderBy("id")
	} else {
		qb.FromRaw("`table_test_from` as `t`").
			Where("email", "like", "%@yao.run").
			Select("id", "t.email as wid", "t.cate as category").
			OrderBy("id")
	}

	// select `id`, `t`.`email` as `wid`, `t`.`cate` as `category` from `table_test_from` as `t` where `email` like ? order by `id` asc
	// select "id", "t"."email" as "wid", "t"."cate" as "category" from "table_test_from" as "t" where "email" like $1 order by "id" asc
	checkingFromFrom(t, qb)
}

func TestFromFromSub(t *testing.T) {
	NewTableFoFromTest()
	qb := getTestBuilder()
	qb.FromSub(func(sub Query) {
		sub.From("table_test_from as t").
			Where("email", "like", "%@yao.run").
			Select("id", "t.email as wid", "t.cate as category").
			OrderBy("id")
	}, "t").
		Where("t.category", "dog").
		OrderByDesc("t.id")

	// select * from (select `id`, `t`.`email` as `wid`, `t`.`cate` as `category` from `table_test_from` as `t` where `email` like ? order by `id` asc) as `t` where `t`.`category` = ? order by `t`.`id` desc
	// select * from (select "id", "t"."email" as "wid", "t"."cate" as "category" from "table_test_from" as "t" where "email" like $1 order by "id" asc) as "t" where "t"."category" = $2 order by "t"."id" desc

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select * from (select "id", "t"."email" as "wid", "t"."cate" as "category" from "table_test_from" as "t" where "email" like $1 order by "id" asc) as "t" where "t"."category" = $2 order by "t"."id" desc`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select * from (select `id`, `t`.`email` as `wid`, `t`.`cate` as `category` from `table_test_from` as `t` where `email` like ? order by `id` asc) as `t` where `t`.`category` = ? order by `t`.`id` desc", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should be have 4 rows")
	if len(rows) == 2 {
		assert.Equal(t, int64(3), rows[0]["id"].(int64), "the id of first row should be 3")
		assert.Equal(t, int64(2), rows[1]["id"].(int64), "the id of first row should be 2")
	}
}

// clean the test data
func TestFromClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_from")
}

func NewTableFoFromTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_from")
	builder.MustCreateTable("table_test_from", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email")
		table.String("name")
		table.String("cate")
		table.Integer("vote")
		table.Float("score", 5, 2)
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})

	qb := getTestBuilder()
	qb.Table("table_test_from").Insert([]xun.R{
		{"email": "john@yao.run", "cate": "cat", "name": "John", "vote": 8, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "cate": "dog", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "cate": "dog", "name": "Ken", "vote": 5, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "cate": "cat", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}

func checkingFromFrom(t *testing.T, qb Query) {
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "t"."email" as "wid", "t"."cate" as "category" from "table_test_from" as "t" where "email" like $1 order by "id" asc`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `t`.`email` as `wid`, `t`.`cate` as `category` from `table_test_from` as `t` where `email` like ? order by `id` asc", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 4, len(rows), "the return value should be have 4 rows")
	if len(rows) == 4 {
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, "john@yao.run", rows[0]["wid"].(string), "the wid of first row should be john@yao.run")
		assert.Equal(t, "cat", rows[0]["category"].(string), "the category of first row should be 1")
	}
}
