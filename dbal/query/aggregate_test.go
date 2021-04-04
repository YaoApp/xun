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

func TestAggregateMustMin(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		MustMin("score")

	assert.Equal(t, float64(48.12), value.MustToFixed(2), "the return value should be 48.12")
}

func TestAggregateMustMax(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		MustMax("vote")
	assert.Equal(t, 125, value.MustInt(), "the return value should be 125")
}
func TestAggregateMustSum(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		MustSum("score")

	assert.Equal(t, float64(308.27), value.MustToFixed(2), "the return value should be 308.27")
}

func TestAggregateMustAvg(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		MustAvg("score")
	assert.Equal(t, float64(77.07), value.MustToFixed(2), "the return value should be 77.07")
}

func TestAggregateUnionMustCount(t *testing.T) {
	NewTableFoAggregateTest()
	qb := getTestBuilder()
	value := qb.Table("table_test_aggregate_t1").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "name").
		Union(func(qb Query) {
			qb.Table("table_test_aggregate_t2").
				Where("email", "like", "%yaojs.org").
				Select("id", "email", "name")
		}).
		MustCount("id")
	assert.Equal(t, int64(7), value, "the return value should be 7")
}

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
