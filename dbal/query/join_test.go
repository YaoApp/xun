package query

import (
	"fmt"
	"testing"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestJoinColumnIsString(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		Join("table_test_join_t2 as t2", "t2.t1_id", "=", "t1.id").
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.title as title", "t2.content as content")

	fmt.Println(qb.ToSQL())
}

// clean the test data
func TestJoinClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_join")
}

func NewTableFoJoinTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_join_t1")
	builder.MustCreateTable("table_test_join_t1", func(table schema.Blueprint) {
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
	qb.Table("table_test_join_t1").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})

	builder.DropTableIfExists("table_test_join_t2")
	builder.MustCreateTable("table_test_join_t2", func(table schema.Blueprint) {
		table.ID("id")
		table.BigInteger("t1_id")
		table.String("title", 300)
		table.LongText("content")
		table.Enum("status", []string{"PUBLISHED", "DRAFT"}).SetDefault("DRAFT")
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_join_t2").Insert([]xun.R{
		{
			"t1_id": 1,
			"title": "A Psychological Trick to Evoke An Interesting Conversation",
			"content": `
				Imagine you pass by a question that asks “the Titanic got invaded by aliens, right?” 
				One one hand you’re holding in a chuckle and slightly in disbelief; on the other hand, 
				you went through the pain of researching and answering the question into such enormous detail. So, what happened?
			`,
			"status": "PUBLISHED", "created_at": "2021-03-26 00:00:16",
		},
		{
			"t1_id": 1,
			"title": "Three Things in Life That Aren’t Worth The Effort",
			"content": `
				To be more efficient and happy, cut the waste and damaging activities from your life.
			`,
			"status": "DRAFT", "created_at": "2021-03-26 00:08:15",
		},
		{
			"t1_id": 2,
			"title": "I tried planking for 5 minutes every day for a month — here’s what happened",
			"content": `
				My core strength improved and my back felt awesome … but the journey wasn’t easy.
			`,
			"status": "DRAFT", "created_at": "2021-03-26 12:35:12",
		},
		{
			"t1_id": 3,
			"title": "The Future of Dashboards is Dashboardless",
			"content": `
				Earlier this year I wrote a long read about where I feel data visualisation is in 2021. In this post, I mention two things that I want to…
			`,
			"status": "PUBLISHED", "created_at": "2021-03-26 19:22:52",
		},
	})
}
