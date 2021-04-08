package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestJoinJoinOn(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		Join("table_test_join_t2 as t2", "t2.t1_id", "=", "t1.id").
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.t1_id", "t2.title as title", "t2.content as content")

	// select `t1`.*, `t2`.`t1_id`, `t2`.`title` as `title`, `t2`.`content` as `content` from `table_test_join_t1` as `t1` inner join `table_test_join_t2` as `t2` on `t2`.`t1_id` = `t1`.`id` where `t2`.`status` = ?
	// select "t1".*, "t2"."t1_id", "t2"."title" as "title", "t2"."content" as "content" from "table_test_join_t1" as "t1" inner join "table_test_join_t2" as "t2" on "t2"."t1_id" = "t1"."id" where "t2"."status" = $1

	// qb.DD()

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "t1".*, "t2"."t1_id", "t2"."title" as "title", "t2"."content" as "content" from "table_test_join_t1" as "t1" inner join "table_test_join_t2" as "t2" on "t2"."t1_id" = "t1"."id" where "t2"."status" = $1`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `t1`.*, `t2`.`t1_id`, `t2`.`title` as `title`, `t2`.`content` as `content` from `table_test_join_t1` as `t1` inner join `table_test_join_t2` as `t2` on `t2`.`t1_id` = `t1`.`id` where `t2`.`status` = ?", sql, "the query sql not equal")
	}

	bindings := qb.GetBindings()
	assert.Equal(t, 1, len(bindings), "the bindings should have 1 item")
	if len(bindings) == 1 {
		assert.Equal(t, "PUBLISHED", bindings[0].(string), "the 1st binding should be PUBLISHED")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should has 1 row")
	if len(rows) == 2 {
		assert.Equal(t, int64(1), rows[0]["t1_id"].(int64), "the t1_id of first row should be 1")
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, "A Psychological Trick to Evoke An Interesting Conversation", rows[0]["title"].(string), "the title of first row should be A Psychological Trick to Evoke An Interesting Conversation")
		assert.Equal(t, int64(4), rows[1]["t1_id"].(int64), "the t1_id of 2nd row should be 4")
		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the id of 2nd should be 4")
		assert.Equal(t, "The Future of Dashboards is Dashboardless", rows[1]["title"].(string), "the title of 2nd row should be The Future of Dashboards is Dashboardless")
	}

}

func TestJoinJoinOnNested(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		Join("table_test_join_t2 as t2", func(join Query) {
			join.On("t2.t1_id", "t1.id").OrOn("t2.o_id", "=", "t1.id")
		}).
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.t1_id", "t2.title as title", "t2.o_id as o_id").
		OrderByDesc("o_id")

	// qb.DD()

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "t1".*, "t2"."t1_id", "t2"."title" as "title", "t2"."o_id" as "o_id" from "table_test_join_t1" as "t1" inner join "table_test_join_t2" as "t2" on ("t2"."t1_id" = "t1"."id" or "t2"."o_id" = "t1"."id") where "t2"."status" = $1 order by "o_id" desc`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `t1`.*, `t2`.`t1_id`, `t2`.`title` as `title`, `t2`.`o_id` as `o_id` from `table_test_join_t1` as `t1` inner join `table_test_join_t2` as `t2` on (`t2`.`t1_id` = `t1`.`id` or `t2`.`o_id` = `t1`.`id`) where `t2`.`status` = ? order by `o_id` desc", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 3, len(rows), "the return value should have 3 rows")
	if len(rows) == 3 {
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, int64(1), rows[0]["t1_id"].(int64), "the t1_id of first row should be 1")
		assert.Equal(t, int64(9), rows[0]["o_id"].(int64), "the o_id of first row should be 9")

		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the join_id of first row should be 4")
		assert.Equal(t, int64(4), rows[1]["t1_id"].(int64), "the t1_id of first row should be 4")
		assert.Equal(t, int64(6), rows[1]["o_id"].(int64), "the o_id of first row should be 6")

		assert.Equal(t, int64(1), rows[2]["id"].(int64), "the join_id of first row should be 1")
		assert.Equal(t, int64(5), rows[2]["t1_id"].(int64), "the t1_id of first row should be 5")
		assert.Equal(t, int64(1), rows[2]["o_id"].(int64), "the o_id of first row should be 1")
	}
}

func TestJoinJoinSubOn(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		JoinSub(func(qb Query) {
			qb.Select("t1_id as join_id", "title", "status").
				From("table_test_join_t2")
		}, "t2", "t2.join_id", "=", "t1.id").
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.join_id", "t2.title as title")

	// qb.DD()

	// select `t1`.*, `t2`.`join_id`, `t2`.`title` as `title` from `table_test_join_t1` as `t1` inner join (select `t1_id` as `join_id`, `title`, `status` from `table_test_join_t2`) as t2 on `t2`.`join_id` = `t1`.`id` where `t2`.`status` = ?
	// select "t1".*, "t2"."join_id", "t2"."title" as "title" from "table_test_join_t1" as "t1" inner join (select "t1_id" as "join_id", "title", "status" from "table_test_join_t2") as t2 on "t2"."join_id" = "t1"."id" where "t2"."status" = $1

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "t1".*, "t2"."join_id", "t2"."title" as "title" from "table_test_join_t1" as "t1" inner join (select "t1_id" as "join_id", "title", "status" from "table_test_join_t2") as t2 on "t2"."join_id" = "t1"."id" where "t2"."status" = $1`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `t1`.*, `t2`.`join_id`, `t2`.`title` as `title` from `table_test_join_t1` as `t1` inner join (select `t1_id` as `join_id`, `title`, `status` from `table_test_join_t2`) as t2 on `t2`.`join_id` = `t1`.`id` where `t2`.`status` = ?", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should has 1 row")
	if len(rows) == 2 {
		assert.Equal(t, int64(1), rows[0]["join_id"].(int64), "the join_id of first row should be 1")
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, "A Psychological Trick to Evoke An Interesting Conversation", rows[0]["title"].(string), "the title of first row should be A Psychological Trick to Evoke An Interesting Conversation")
		assert.Equal(t, int64(4), rows[1]["join_id"].(int64), "the join_id of 2nd row should be 4")
		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the id of 2nd should be 4")
		assert.Equal(t, "The Future of Dashboards is Dashboardless", rows[1]["title"].(string), "the title of 2nd row should be The Future of Dashboards is Dashboardless")
	}

}

func TestJoinLeftJoinOn(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		LeftJoin("table_test_join_t2 as t2", "t2.t1_id", "=", "t1.id").
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.t1_id", "t2.title as title", "t2.content as content")

	// select `t1`.*, `t2`.`t1_id`, `t2`.`title` as `title`, `t2`.`content` as `content` from `table_test_join_t1` as `t1` left join `table_test_join_t2` as `t2` on `t2`.`t1_id` = `t1`.`id` where `t2`.`status` = ?
	// select "t1".*, "t2"."t1_id", "t2"."title" as "title", "t2"."content" as "content" from "table_test_join_t1" as "t1" left join "table_test_join_t2" as "t2" on "t2"."t1_id" = "t1"."id" where "t2"."status" = $1

	// qb.DD()

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "t1".*, "t2"."t1_id", "t2"."title" as "title", "t2"."content" as "content" from "table_test_join_t1" as "t1" left join "table_test_join_t2" as "t2" on "t2"."t1_id" = "t1"."id" where "t2"."status" = $1`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `t1`.*, `t2`.`t1_id`, `t2`.`title` as `title`, `t2`.`content` as `content` from `table_test_join_t1` as `t1` left join `table_test_join_t2` as `t2` on `t2`.`t1_id` = `t1`.`id` where `t2`.`status` = ?", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should has 1 row")
	if len(rows) == 2 {
		assert.Equal(t, int64(1), rows[0]["t1_id"].(int64), "the t1_id of first row should be 1")
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, "A Psychological Trick to Evoke An Interesting Conversation", rows[0]["title"].(string), "the title of first row should be A Psychological Trick to Evoke An Interesting Conversation")
		assert.Equal(t, int64(4), rows[1]["t1_id"].(int64), "the t1_id of 2nd row should be 4")
		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the id of 2nd should be 4")
		assert.Equal(t, "The Future of Dashboards is Dashboardless", rows[1]["title"].(string), "the title of 2nd row should be The Future of Dashboards is Dashboardless")
	}

}

func TestJoinLeftJoinSubOn(t *testing.T) {
	NewTableFoJoinTest()
	qb := getTestBuilder()
	qb.Table("table_test_join_t1 as t1").
		LeftJoinSub(func(qb Query) {
			qb.Select("t1_id as join_id", "title", "status").
				From("table_test_join_t2")
		}, "t2", "t2.join_id", "=", "t1.id").
		Where("t2.status", "=", "PUBLISHED").
		Select("t1.*", "t2.join_id", "t2.title as title")

	// qb.DD()

	// select `t1`.*, `t2`.`join_id`, `t2`.`title` as `title` from `table_test_join_t1` as `t1` left join (select `t1_id` as `join_id`, `title`, `status` from `table_test_join_t2`) as t2 on `t2`.`join_id` = `t1`.`id` where `t2`.`status` = ?
	// select "t1".*, "t2"."join_id", "t2"."title" as "title" from "table_test_join_t1" as "t1" left join (select "t1_id" as "join_id", "title", "status" from "table_test_join_t2") as t2 on "t2"."join_id" = "t1"."id" where "t2"."status" = $1

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select "t1".*, "t2"."join_id", "t2"."title" as "title" from "table_test_join_t1" as "t1" left join (select "t1_id" as "join_id", "title", "status" from "table_test_join_t2") as t2 on "t2"."join_id" = "t1"."id" where "t2"."status" = $1`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select `t1`.*, `t2`.`join_id`, `t2`.`title` as `title` from `table_test_join_t1` as `t1` left join (select `t1_id` as `join_id`, `title`, `status` from `table_test_join_t2`) as t2 on `t2`.`join_id` = `t1`.`id` where `t2`.`status` = ?", sql, "the query sql not equal")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 2, len(rows), "the return value should has 1 row")
	if len(rows) == 2 {
		assert.Equal(t, int64(1), rows[0]["join_id"].(int64), "the join_id of first row should be 1")
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the id of first row should be 1")
		assert.Equal(t, "A Psychological Trick to Evoke An Interesting Conversation", rows[0]["title"].(string), "the title of first row should be A Psychological Trick to Evoke An Interesting Conversation")
		assert.Equal(t, int64(4), rows[1]["join_id"].(int64), "the join_id of 2nd row should be 4")
		assert.Equal(t, int64(4), rows[1]["id"].(int64), "the id of 2nd should be 4")
		assert.Equal(t, "The Future of Dashboards is Dashboardless", rows[1]["title"].(string), "the title of 2nd row should be The Future of Dashboards is Dashboardless")
	}

}

// clean the test data
func TestJoinClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_join_t1")
	builder.DropTableIfExists("table_test_join_t2")
	builder.DropTableIfExists("table_test_join_t3")
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
		table.BigInteger("o_id")
		table.String("title", 300)
		table.LongText("content")
		table.Enum("status", []string{"PUBLISHED", "DRAFT"}).SetDefault("DRAFT")
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_join_t2").Insert([]xun.R{
		{
			"t1_id": 1,
			"o_id":  9,
			"title": "A Psychological Trick to Evoke An Interesting Conversation",
			"content": `
				Imagine you pass by a question that asks “the Titanic got invaded by aliens, right?
				One one hand you’re holding in a chuckle and slightly in disbelief; on the other hand, 
				you went through the pain of researching and answering the question into such enormous detail. So, what happened?
			`,
			"status": "PUBLISHED", "created_at": "2021-03-26 00:00:16",
		},
		{
			"t1_id": 2,
			"o_id":  8,
			"title": "Three Things in Life That Aren’t Worth The Effort",
			"content": `
				To be more efficient and happy, cut the waste and damaging activities from your life.
			`,
			"status": "DRAFT", "created_at": "2021-03-26 00:08:15",
		},
		{
			"t1_id": 3,
			"o_id":  7,
			"title": "I tried planking for 5 minutes every day for a month — here’s what happened",
			"content": `
				My core strength improved and my back felt awesome … but the journey wasn’t easy.
			`,
			"status": "DRAFT", "created_at": "2021-03-26 12:35:12",
		},
		{
			"t1_id": 4,
			"o_id":  6,
			"title": "The Future of Dashboards is Dashboardless",
			"content": `
				Earlier this year I wrote a long read about where I feel data visualisation is in 2021. In this post, I mention two things that I want to…
			`,
			"status": "PUBLISHED", "created_at": "2021-03-26 19:22:52",
		}, {
			"t1_id": 5,
			"o_id":  1,
			"title": "The Future of Dashboards is Dashboardless",
			"content": `
				Earlier this year I wrote a long read about where I feel data visualisation is in 2021. In this post, I mention two things that I want to…
			`,
			"status": "PUBLISHED", "created_at": "2021-03-26 19:22:52",
		},
	})

	builder.DropTableIfExists("table_test_join_t3")
	builder.MustCreateTable("table_test_join_t3", func(table schema.Blueprint) {
		table.ID("id")
		table.BigInteger("t2_id")
		table.String("Name", 300)
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_join_t3").Insert([]xun.R{
		{"t2_id": 1, "Name": "Emma", "created_at": "2021-03-27 00:00:16"},
		{"t2_id": 2, "Name": "Ava", "created_at": "2021-03-27 08:13:23"},
		{"t2_id": 3, "Name": "Amelia", "created_at": "2021-03-27 09:46:21"},
		{"t2_id": 4, "Name": "Elizabeth", "created_at": "2021-03-27 14:00:22"},
	})
}
