package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestGroupGroupBy(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	qb.Table("table_test_group").
		Where("email", "like", "%@yao.run").
		Select("status", dbal.Raw("Count(id) as cnt")).
		GroupBy("status")

	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status`
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status"
	checkGroupGroupBy(t, qb)
}

func TestGroupGroupByRaw(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupByRaw(`"status"`)
	} else {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupByRaw("`status`")
	}

	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status`
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status"
	checkGroupGroupBy(t, qb)
}

func TestGroupHaving(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	qb.Table("table_test_group").
		Where("email", "like", "%@yao.run").
		Select("status", dbal.Raw("Count(id) as cnt")).
		GroupBy("status").
		Having("status", "=", "DONE")

	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ?
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2
	checkGroupHaving(t, qb)
}

func TestGroupOrHaving(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	qb.Table("table_test_group").
		Where("email", "like", "%@yao.run").
		Select("status", dbal.Raw("Count(id) as cnt")).
		GroupBy("status").
		Having("status", "=", "DONE").
		OrHaving("status", "=", "PENDING")

	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ? or `status` = ?
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2 or "status" = $3
	checkGroupOrHaving(t, qb)
}

func TestGroupHavingRaw(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupBy("status").
			HavingRaw(`"status" = $2`, "DONE")
	} else {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupBy("status").
			HavingRaw("`status` = ?", "DONE")
	}
	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ?
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2
	checkGroupHaving(t, qb)
}

func TestGroupOrHavingRaw(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupBy("status").
			Having("status", "=", "DONE").
			OrHavingRaw(`"status" = $3`, "PENDING")
	} else {
		qb.Table("table_test_group").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupBy("status").
			Having("status", "=", "DONE").
			OrHavingRaw("`status` = ?", "PENDING")
	}

	// select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ? or `status` = ?
	// select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2 or "status" = $3
	checkGroupOrHaving(t, qb)
}

func TestGroupHavingBetween(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	qb.Table("table_test_group").
		Where("email", "like", "%@yao.run").
		Select("vote", dbal.Raw("Count(id) as cnt")).
		GroupBy("vote").
		HavingBetween("vote", []interface{}{5, 7})

	// select `vote`, Count(id) as cnt from `table_test_group` where `email` like ? group by `vote` having `vote` between ? and ?
	// select "vote", Count(id) as cnt from "table_test_group" where "email" like $1 group by "vote" having "vote" between $2 and $3

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "vote", Count(id) as cnt from "table_test_group" where "email" like $1 group by "vote" having "vote" between $2 and $3`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `vote`, Count(id) as cnt from `table_test_group` where `email` like ? group by `vote` having `vote` between ? and ?", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should have 2 items")
	if len(rows) == 2 {
		assert.Equal(t, int64(2), rows[0]["cnt"].(int64), "the cnt of first row should be 1")
		assert.Equal(t, int64(5), rows[0]["vote"].(int64), "the vote of first row should be 5")
		assert.Equal(t, int64(1), rows[1]["cnt"].(int64), "the cnt of last row should be 2")
		assert.Equal(t, int64(6), rows[1]["vote"].(int64), "the vote of last row should be DONE")
	}
}

func TestGroupOrHavingBetween(t *testing.T) {
	NewTableFoGroupTest()
	qb := getTestBuilder()
	qb.Table("table_test_group").
		Where("email", "like", "%@yao.run").
		Select("vote", dbal.Raw("Count(id) as cnt")).
		GroupBy("vote").
		Having("vote", "=", 8).
		OrHavingBetween("vote", []interface{}{5, 7})

	// select `vote`, Count(id) as cnt from `table_test_group` where `email` like ? group by `vote` having `vote` = ? or `vote` between ? and ?
	// select "vote", Count(id) as cnt from "table_test_group" where "email" like $1 group by "vote" having "vote" = $2 or "vote" between $3 and $4

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "vote", Count(id) as cnt from "table_test_group" where "email" like $1 group by "vote" having "vote" = $2 or "vote" between $3 and $4`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `vote`, Count(id) as cnt from `table_test_group` where `email` like ? group by `vote` having `vote` = ? or `vote` between ? and ?", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 3, len(rows), "the return value should have 2 items")
	if len(rows) == 3 {
		assert.Equal(t, int64(2), rows[0]["cnt"].(int64), "the cnt of first item should be 1")
		assert.Equal(t, int64(5), rows[0]["vote"].(int64), "the vote of first item should be 5")
		assert.Equal(t, int64(1), rows[1]["cnt"].(int64), "the cnt of second item should be 2")
		assert.Equal(t, int64(6), rows[1]["vote"].(int64), "the vote of second item should be DONE")
		assert.Equal(t, int64(1), rows[2]["cnt"].(int64), "the cnt of last item should be 2")
		assert.Equal(t, int64(8), rows[2]["vote"].(int64), "the vote of last item should be DONE")
	}
}

// clean the test data
func TestGroupClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_group")
}

func NewTableFoGroupTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_group")
	builder.MustCreateTable("table_test_group", func(table schema.Blueprint) {
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
	qb.Table("table_test_group").Insert([]xun.R{
		{"email": "john@yao.run", "cate": "cat", "name": "John", "vote": 8, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "cate": "dog", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "cate": "dog", "name": "Ken", "vote": 5, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "cate": "cat", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}

func checkGroupHaving(t *testing.T, qb Query) {

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ?", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 1, len(rows), "the return value should has 1 row")
	if len(rows) == 1 {
		assert.Equal(t, int64(2), rows[0]["cnt"].(int64), "the cnt of first row should be 2")
		assert.Equal(t, "DONE", rows[0]["status"].(string), "the status of first row should be DONE")
	}

}

func checkGroupOrHaving(t *testing.T, qb Query) {
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status" having "status" = $2 or "status" = $3`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status` having `status` = ? or `status` = ?", sql, "the query sql not equal")
	}

	// // check values
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should has 1 row")
	if len(rows) == 2 {
		if unit.DriverIs("sqlite3") {
			assert.Equal(t, int64(2), rows[0]["cnt"].(int64), "the cnt of first row should be 2")
			assert.Equal(t, "DONE", rows[0]["status"].(string), "the status of first row should be DONE")
			assert.Equal(t, int64(1), rows[1]["cnt"].(int64), "the cnt of last row should be 1")
			assert.Equal(t, "PENDING", rows[1]["status"].(string), "the status of last row should be PENDING")
		} else {
			assert.Equal(t, int64(1), rows[0]["cnt"].(int64), "the cnt of first row should be 1")
			assert.Equal(t, "PENDING", rows[0]["status"].(string), "the status of first row should be PENDING")
			assert.Equal(t, int64(2), rows[1]["cnt"].(int64), "the cnt of last row should be 2")
			assert.Equal(t, "DONE", rows[1]["status"].(string), "the status of last row should be DONE")
		}
	}
}

func checkGroupGroupBy(t *testing.T, qb Query) {
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "status", Count(id) as cnt from "table_test_group" where "email" like $1 group by "status"`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `status`, Count(id) as cnt from `table_test_group` where `email` like ? group by `status`", sql, "the query sql not equal")
	}

	// check values
	rows := qb.MustGet()
	assert.Equal(t, 3, len(rows), "the return value should has 3 rows")
	if len(rows) == 3 {
		if unit.DriverIs("sqlite3") {
			assert.Equal(t, int64(2), rows[0]["cnt"].(int64), "the cnt of first row should be 2")
			assert.Equal(t, "DONE", rows[0]["status"].(string), "the status of first row should be DONE")
		} else {
			assert.Equal(t, int64(2), rows[2]["cnt"].(int64), "the cnt of last row should be 2")
			assert.Equal(t, "DONE", rows[2]["status"].(string), "the status of last row should be DONE")
		}
	}
}
