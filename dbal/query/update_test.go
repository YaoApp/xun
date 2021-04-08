package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestInsertMustUpsert(t *testing.T) {
	NewTableForUpdateTest()
	qb := getTestBuilder()
	affected := qb.Table("table_test_update").MustUpsert([]xun.R{
		{"email": "max@yao.run", "name": "Max", "vote": 19, "score": 86.32, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
		{"email": "john@yao.run", "name": "John", "vote": 20, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
	}, []string{"email"}, []string{"vote"})

	if unit.DriverIs("mysql") {
		assert.Equal(t, int64(3), affected, "The affected rows should be 3")
	} else if unit.DriverIs("postgres") {
		assert.Equal(t, int64(2), affected, "The affected rows should be 2")
	}
}

func TestInsertMustUpsertUpdateValue(t *testing.T) {
	NewTableForUpdateTest()
	qb := getTestBuilder()
	affected := qb.Table("table_test_update").MustUpsert([]xun.R{
		{"email": "max@yao.run", "name": "Max", "vote": 19, "score": 86.32, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
		{"email": "john@yao.run", "name": "John", "vote": 20, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
	}, []string{"email"}, map[string]interface{}{"vote": 100, "score": 99.98})

	if unit.DriverIs("mysql") {
		assert.Equal(t, int64(3), affected, "The affected rows should be 3")
	} else if unit.DriverIs("postgres") {
		assert.Equal(t, int64(2), affected, "The affected rows should be 2")
	}
}

func TestInsertMustUpsertUpdateRaw(t *testing.T) {
	NewTableForUpdateTest()
	qb := getTestBuilder()

	raw := "random()"
	if unit.DriverIs("mysql") {
		raw = "rand()"
	}

	affected := qb.Table("table_test_update").MustUpsert([]xun.R{
		{"email": "max@yao.run", "name": "Max", "vote": 19, "score": 86.32, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
		{"email": "john@yao.run", "name": "John", "vote": 20, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-27 07:16:16", "updated_at": "2021-03-27 07:16:16"},
	}, []string{"email"}, map[string]interface{}{"vote": dbal.Raw(raw), "score": 99.98})

	if unit.DriverIs("mysql") {
		assert.Equal(t, int64(3), affected, "The affected rows should be 3")
	} else if unit.DriverIs("postgres") {
		assert.Equal(t, int64(2), affected, "The affected rows should be 2")
	}
}

// clean the test data
func TestUpdateClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_update")
}

func NewTableForUpdateTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_update")
	builder.MustCreateTable("table_test_update", func(table schema.Blueprint) {
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
	qb.Table("table_test_update").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
