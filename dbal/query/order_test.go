package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestOrderOrderBy(t *testing.T) {
	NewTableForOrderTest()
	qb := getTestBuilder()
	qb.Table("table_test_order").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email", "vote", "score", "status").
		OrderBy("vote", "desc").
		OrderBy("score")

	// select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? order by `vote` desc, `score` asc
	// select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 order by "vote" desc, "score" asc
	checkOrderOrderBy(t, qb)
}

func TestOrderOrderByDesc(t *testing.T) {
	NewTableForOrderTest()
	qb := getTestBuilder()
	qb.Table("table_test_order").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email", "vote", "score", "status").
		OrderByDesc("vote").
		OrderBy("score")

	// select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? order by `vote` desc, `score` asc
	// select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 order by "vote" desc, "score" asc
	checkOrderOrderBy(t, qb)
}

func TestOrderOrderByRaw(t *testing.T) {
	NewTableForOrderTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_order").
			Where("email", "like", "%@yao.run").
			Select("id", "name", "email", "vote", "score", "status").
			OrderByRaw(`"vote" desc, "score" asc`)
	} else {
		qb.Table("table_test_order").
			Where("email", "like", "%@yao.run").
			Select("id", "name", "email", "vote", "score", "status").
			OrderByRaw("`vote` desc, `score` asc")
	}

	// select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? order by `vote` ?, `score` asc
	// select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 order by "vote" $2, "score" asc
	checkOrderOrderBy(t, qb)
}

func TestOrderOrderByUnion(t *testing.T) {
	NewTableForOrderTest()
	qb := getTestBuilder()
	qb.Table("table_test_order").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email", "vote", "score", "status").
		Where("vote", 5).
		Union(func(qb Query) {
			qb.Table("table_test_order").
				Select("id", "name", "email", "vote", "score", "status").
				Where("vote", 6)
		}).
		OrderBy("vote", "desc").
		OrderBy("score")

	// (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc
	// (select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 and "vote" = $2 ) union (select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "vote" = $3) order by "vote" desc, "score" asc
	// select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc
	checkOrderOrderByUnion(t, qb)
}

func TestOrderOrderByRawUnion(t *testing.T) {
	NewTableForOrderTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_order").
			Where("email", "like", "%@yao.run").
			Select("id", "name", "email", "vote", "score", "status").
			Where("vote", 5).
			Union(func(qb Query) {
				qb.Table("table_test_order").
					Select("id", "name", "email", "vote", "score", "status").
					Where("vote", 6)
			}).
			OrderByRaw(`"vote" desc, "score" asc`)
	} else {
		qb.Table("table_test_order").
			Where("email", "like", "%@yao.run").
			Select("id", "name", "email", "vote", "score", "status").
			Where("vote", 5).
			Union(func(qb Query) {
				qb.Table("table_test_order").
					Select("id", "name", "email", "vote", "score", "status").
					Where("vote", 6)
			}).
			OrderByRaw("`vote` desc, `score` asc")
	}
	// (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc
	// (select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 and "vote" = $2 ) union (select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "vote" = $3) order by "vote" desc, "score" asc
	// select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc
	checkOrderOrderByUnion(t, qb)
}

// clean the test data
func TestOrderClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_order")
}

func NewTableForOrderTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_order")
	builder.MustCreateTable("table_test_order", func(table schema.Blueprint) {
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
	qb.Table("table_test_order").Insert([]xun.R{
		{"email": "john@yao.run", "cate": "cat", "name": "John", "vote": 8, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "cate": "dog", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "cate": "dog", "name": "Ken", "vote": 5, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "cate": "cat", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}

func checkOrderOrderBy(t *testing.T, qb Query) {
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 order by "vote" desc, "score" asc`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? order by `vote` desc, `score` asc", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 4, len(rows), "the return value should be have 4 rows")
	if len(rows) == 4 {
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of 1st row should be 1")
		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the id of 2nd row should be 4")
		assert.Equal(t, int64(2), rows[2]["id"].(int64), "the id of 3rd row should be 2")
		assert.Equal(t, int64(3), rows[3]["id"].(int64), "the id of 4th row should be 3")
	}
}

func checkOrderOrderByUnion(t *testing.T, qb Query) {
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `(select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "email" like $1 and "vote" = $2 ) union (select "id", "name", "email", "vote", "score", "status" from "table_test_order" where "vote" = $3) order by "vote" desc, "score" asc`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union select * from (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "(select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `email` like ? and `vote` = ? ) union (select `id`, `name`, `email`, `vote`, `score`, `status` from `table_test_order` where `vote` = ?) order by `vote` desc, `score` asc", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 3, len(rows), "the return value should be have 3 rows")
	if len(rows) == 3 {
		assert.Equal(t, int64(4), rows[0]["id"].(int64), "the id of 1st row should be 4")
		assert.Equal(t, int64(2), rows[1]["id"].(int64), "the id of 2nd row should be 2")
		assert.Equal(t, int64(3), rows[2]["id"].(int64), "the id of 3rd row should be 3")
	}
}
