package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestLimitTake(t *testing.T) {
	NewTableForLimitTest()
	qb := getTestBuilder()
	qb.Table("table_test_limit").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email").
		OrderByDesc("id").
		Take(2)

	// select `id`, `name`, `email` from `table_test_limit` where `email` like ? order by `id` desc limit 2
	// select "id", "name", "email" from "table_test_limit" where "email" like $1 order by "id" asc limit 2

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "name", "email" from "table_test_limit" where "email" like $1 order by "id" desc limit 2`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `name`, `email` from `table_test_limit` where `email` like ? order by `id` desc limit 2", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should be have 2 rows")
	if len(rows) == 2 {
		assert.Equal(t, int64(4), rows[0]["id"].(int64), "the id of 1st row should be 4")
		assert.Equal(t, int64(3), rows[1]["id"].(int64), "the id of 2nd row should be 3")
	}
}

func TestLimitTakeUnion(t *testing.T) {
	NewTableForLimitTest()
	qb := getTestBuilder()
	qb.Table("table_test_limit").
		Select("id", "name", "email").
		Where("id", 2).
		Union(func(qb Query) {
			qb.Table("table_test_limit").
				Select("id", "name", "email").
				Where("id", 4)
		}).
		OrderByDesc("id").
		Take(1)

	// (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1
	// (select "id", "name", "email" from "table_test_limit" where "id" = $1 ) union (select "id", "name", "email" from "table_test_limit" where "id" = $2) order by "id" desc limit 1
	// select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `(select "id", "name", "email" from "table_test_limit" where "id" = $1 ) union (select "id", "name", "email" from "table_test_limit" where "id" = $2) order by "id" desc limit 1`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "(select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 1, len(rows), "the return value should be have 2 rows")
	if len(rows) == 1 {
		assert.Equal(t, int64(4), rows[0]["id"].(int64), "the id of 1st row should be 4")
	}
}

func TestLimitSkip(t *testing.T) {
	NewTableForLimitTest()
	qb := getTestBuilder()
	qb.Table("table_test_limit").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email").
		OrderByDesc("id").
		Skip(1).Take(2)

	// select `id`, `name`, `email` from `table_test_limit` where `email` like ? order by `id` desc limit 2 offset 1
	// select "id", "name", "email" from "table_test_limit" where "email" like $1 order by "id" asc limit 2 offset 1

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "name", "email" from "table_test_limit" where "email" like $1 order by "id" desc limit 2 offset 1`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `name`, `email` from `table_test_limit` where `email` like ? order by `id` desc limit 2 offset 1", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should be have 2 rows")
	if len(rows) == 2 {
		assert.Equal(t, int64(3), rows[0]["id"].(int64), "the id of 1st row should be 3")
		assert.Equal(t, int64(2), rows[1]["id"].(int64), "the id of 2nd row should be 2")
	}
}

func TestLimitSkipUnion(t *testing.T) {
	NewTableForLimitTest()
	qb := getTestBuilder()
	qb.Table("table_test_limit").
		Select("id", "name", "email").
		Where("id", 2).
		Union(func(qb Query) {
			qb.Table("table_test_limit").
				Select("id", "name", "email").
				Where("id", 4)
		}).
		OrderByDesc("id").
		Skip(1).Take(1)

	// (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1 offset 1
	// (select "id", "name", "email" from "table_test_limit" where "id" = $1 ) union (select "id", "name", "email" from "table_test_limit" where "id" = $2) order by "id" desc limit 1 offset 1
	// select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1 offset 1

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `(select "id", "name", "email" from "table_test_limit" where "id" = $1 ) union (select "id", "name", "email" from "table_test_limit" where "id" = $2) order by "id" desc limit 1 offset 1`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union select * from (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1 offset 1", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "(select `id`, `name`, `email` from `table_test_limit` where `id` = ? ) union (select `id`, `name`, `email` from `table_test_limit` where `id` = ?) order by `id` desc limit 1 offset 1", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 1, len(rows), "the return value should be have 2 rows")
	if len(rows) == 1 {
		assert.Equal(t, int64(2), rows[0]["id"].(int64), "the id of 1st row should be 2")
	}
}

// clean the test data
func TestLimitClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_limit")
}

func NewTableForLimitTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_limit")
	builder.MustCreateTable("table_test_limit", func(table schema.Blueprint) {
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
	qb.Table("table_test_limit").Insert([]xun.R{
		{"email": "john@yao.run", "cate": "cat", "name": "John", "vote": 8, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "cate": "dog", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "cate": "dog", "name": "Ken", "vote": 5, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "cate": "cat", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
