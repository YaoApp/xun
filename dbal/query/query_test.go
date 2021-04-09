package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestQueryMustExistsTrue(t *testing.T) {
	NewTableFoSelectTest()
	qb := getTestBuilder()
	res := qb.From("table_test_select as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustExists()

	assert.True(t, res, "the return value should be true")
}

func TestQueryMustExistsFalse(t *testing.T) {
	NewTableFoSelectTest()
	qb := getTestBuilder()
	res := qb.From("table_test_select as t").
		Where("email", "like", "%@iqka.com").
		OrderBy("id").
		MustExists()

	assert.False(t, res, "the return value should be false")
}

func TestQueryMustDoesntExistTrue(t *testing.T) {
	NewTableFoSelectTest()
	qb := getTestBuilder()
	res := qb.From("table_test_select as t").
		Where("email", "like", "%@iqka.com").
		OrderBy("id").
		MustDoesntExist()

	assert.True(t, res, "the return value should be true")
}

func TestQueryMustDoesntExistFalse(t *testing.T) {
	NewTableFoSelectTest()
	qb := getTestBuilder()
	res := qb.From("table_test_select as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustDoesntExist()

	assert.False(t, res, "the return value should be false")
}

// clean the test data
func TestQueryClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_query")
}

func NewTableForQueryTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_query")
	builder.MustCreateTable("table_test_query", func(table schema.Blueprint) {
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
	qb.Table("table_test_query").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
