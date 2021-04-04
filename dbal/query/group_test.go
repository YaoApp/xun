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
	NewTableFoUnionTest()
	qb := getTestBuilder()
	qb.Table("table_test_union_t1").
		Where("email", "like", "%@yao.run").
		Select("status", dbal.Raw("Count(id) as cnt")).
		GroupBy("status")

	// select `status`, Count(id) as cnt from `table_test_union_t1` where `email` like ? group by `status`
	// select "status", Count(id) as cnt from "table_test_union_t1" where "email" like $1 group by "status"

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "status", Count(id) as cnt from "table_test_union_t1" where "email" like $1 group by "status"`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
	} else {
		assert.Equal(t, "select `status`, Count(id) as cnt from `table_test_union_t1` where `email` like ? group by `status`", sql, "the query sql not equal")
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

func TestGroupGroupByRaw(t *testing.T) {
	NewTableFoUnionTest()
	qb := getTestBuilder()
	if unit.DriverIs("postgres") {
		qb.Table("table_test_union_t1").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupByRaw(`"status"`)
	} else {
		qb.Table("table_test_union_t1").
			Where("email", "like", "%@yao.run").
			Select("status", dbal.Raw("Count(id) as cnt")).
			GroupByRaw("`status`")
	}

	// select `status`, Count(id) as cnt from `table_test_union_t1` where `email` like ? group by `status`
	// select "status", Count(id) as cnt from "table_test_union_t1" where "email" like $1 group by "status"

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "status", Count(id) as cnt from "table_test_union_t1" where "email" like $1 group by "status"`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
	} else {
		assert.Equal(t, "select `status`, Count(id) as cnt from `table_test_union_t1` where `email` like ? group by `status`", sql, "the query sql not equal")
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
		{"email": "john@yao.run", "cate": "cat", "name": "John", "vote": 10, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "cate": "dog", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "cate": "dog", "name": "Ken", "vote": 125, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "cate": "cat", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
