package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestLockSharedLock(t *testing.T) {
	NewTableForLockTest()
	qb := getTestBuilder()
	qb.UseRead().
		Table("table_test_lock").
		Select("id", "vote").
		OrderByDesc("id")

	assert.True(t, qb.IsRead(), "the connection should be read")

	qb.SharedLock()
	assert.True(t, qb.IsWrite(), "the connection should be write")

	// qb.DD()
	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "vote" from "table_test_lock" order by "id" desc for share`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select `id`, `vote` from `table_test_lock` order by `id` desc", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `vote` from `table_test_lock` order by `id` desc lock in share mode", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 4, len(rows), "the return value should be have 4 items")
}

func TestLockLockForUpdate(t *testing.T) {
	NewTableForLockTest()
	qb := getTestBuilder()
	qb.UseRead().
		Table("table_test_lock").
		Select("id", "vote").
		OrderByDesc("id")

	assert.True(t, qb.IsRead(), "the connection should be read")

	qb.LockForUpdate()
	assert.True(t, qb.IsWrite(), "the connection should be write")

	// qb.DD()

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "id", "vote" from "table_test_lock" order by "id" desc for update`, sql, "the query sql not equal")
	} else if unit.DriverIs("sqlite3") {
		assert.Equal(t, "select `id`, `vote` from `table_test_lock` order by `id` desc", sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `id`, `vote` from `table_test_lock` order by `id` desc for update", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 4, len(rows), "the return value should be have 4 items")
}

// clean the test data
func TestLockClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_lock")
}

func NewTableForLockTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_lock")
	builder.MustCreateTable("table_test_lock", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email").Unique()
		table.String("name").Index()
		table.Integer("vote")
		table.Float("score", 5, 2).Index()
		table.Float("score_grade", 5, 2).Index()
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})

	qb := getTestBuilder()
	qb.Table("table_test_lock").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
