package query

import (
	"testing"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestWhereColumnIsArray(t *testing.T) {
	NewTableFoWhereTest()
	qb := getTestBuilder()
	qb.Table("table_test_where").
		Where("email", "like", "%@yao.run").
		Where([][]interface{}{
			{"score", ">", 64.56},
			{"vote", 10},
		})

	// select `*` from `table_test_where` where `email` like ? and (`score` > ? and `vote` = ?)
	qb.Get()
}

func TestWhereColumnIsClosure(t *testing.T) {
	NewTableFoWhereTest()
	qb := getTestBuilder()
	qb.Table("table_test_where").
		Where("email", "like", "%@yao.run").
		Where(func(qb Query) {
			qb.Where("vote", ">", 10)
			qb.Where("name", "Ken")
			qb.Where(func(qb Query) {
				qb.Where("created_at", ">", "2021-03-25 08:00:00")
				qb.Where("created_at", "<", "2021-03-25 19:00:00")
			})
		}).
		Where("score", ">", 5)

	qb.Get()
	// AND  `email` LIKE '%@yao.run' AND `score` > 5 AND ( `vote` > 10 AND `name` = 'Ken'  AND (`created_at` > '2021-03-25 08:00:00' AND `created_at` < '2021-03-25 19:00:00' ) )
}
func TestWhereValueIsClosure(t *testing.T) {
	NewTableFoWhereTest()
	qb := getTestBuilder()
	qb.Table("table_test_where").
		Where("email", "like", "%@yao.run").
		Where("vote", '>', func(sub Query) {
			sub.From("table_test_where").
				Where("score", ">", 5)
			//   Sum()
		})

	qb.Get()
	// AND  `email` LIKE '%@yao.run' AND `score` > 5 AND ( `vote` > 10 AND `name` = 'Ken'  AND (`created_at` > '2021-03-25 08:00:00' AND `created_at` < '2021-03-25 19:00:00' ) )
}

// clean the test data
func TestWhereClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_where")
}

func NewTableFoWhereTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_where")
	builder.MustCreateTable("table_test_where", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email").Unique()
		table.String("name").Index()
		table.Integer("vote")
		table.Float("score", 5, 2).Index()
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})

	qb := getTestBuilder()
	qb.Table("table_test_where").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
