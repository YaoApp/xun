package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestDeleteMustDelete(t *testing.T) {
	NewTableForDeleteTest()
	qb := getTestBuilder()
	affected := qb.From("table_test_delete").
		Where("id", ">", 2).
		MustDelete()

	assert.Equal(t, int64(2), affected, "The affected rows should be 2")
}

func TestDeleteMustDeleteError(t *testing.T) {
	NewTableForDeleteTest()
	assert.Panics(t, func() {
		newQuery := New(unit.Driver(), unit.DSN())
		newQuery.DB().Close()
		newQuery.From("table_test_delete").
			Where("id", ">", 2).
			MustDelete()
	})
}

func TestDeleteMustDeleteWithJoin(t *testing.T) {
	NewTableForDeleteTest()
	qb := getTestBuilder()
	affected := qb.From("table_test_delete as t1").
		JoinSub(func(qb Query) {
			qb.From("table_test_delete").
				Where("id", ">", 1).
				Select("id as join_id", "score_grade as join_score")
		}, "t2", "t2.join_id", "=", "t1.id").
		Where("t1.id", ">", 2).
		MustDelete()

	assert.Equal(t, int64(2), affected, "The affected rows should be 2")
}

// clean the test data
func TestDeleteClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_delete")
}

func NewTableForDeleteTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_delete")
	builder.MustCreateTable("table_test_delete", func(table schema.Blueprint) {
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
	qb.Table("table_test_delete").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
