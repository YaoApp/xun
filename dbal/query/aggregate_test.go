package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestAggregateMustCount(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		MustCount("id")
	assert.Equal(t, int64(4), value, "the return value should be 4")
}

// @todo: test union

// @todo: test unionOrders

// @todo: test unionLimit

// @todo: test unionOffset

// clean the test data
func TestAggregateClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_aggregate_t1")
	builder.DropTableIfExists("table_test_aggregate_t2")
}

func NewTableFoAggregateTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_aggregate_t1")
	builder.MustCreateTable("table_test_aggregate_t1", func(table schema.Blueprint) {
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
	qb.Table("table_test_aggregate_t1").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})

	builder.DropTableIfExists("table_test_aggregate_t2")
	builder.MustCreateTable("table_test_aggregate_t2", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email")
		table.String("name")
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_aggregate_t2").Insert([]xun.R{
		{"email": "nio@yaojs.org", "name": "Nio", "status": "WAITING", "created_at": "2021-03-26 00:15:16"},
		{"email": "Tom@yaojs.org", "name": "Tom", "status": "PENDING", "created_at": "2021-03-26 08:19:15"},
		{"email": "Han@yaojs.org", "name": "Han", "status": "DONE", "created_at": "2021-03-26 10:24:23"},
	})
}
