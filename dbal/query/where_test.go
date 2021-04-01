package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	//select * from `table_test_where` where `email` like ? and (`score` > ? and `vote` = ?)
	//select * from "table_test_where" where "email" like $1 and ("score" > $2 and "vote" = $3)

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select * from "table_test_where" where "email" like $1 and ("score" > $2 and "vote" = $3)`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select * from `table_test_where` where `email` like ? and (`score` > ? and `vote` = ?)", sql, "the query sql not equal")
	}

	// checking bindings
	bindings := qb.GetBindings()
	assert.Equal(t, 3, len(bindings), "the bindings should have 3 items")
	if len(bindings) == 3 {
		assert.Equal(t, "%@yao.run", bindings[0].(string), "the 1st binding should be %@yao.run")
		assert.Equal(t, float64(64.56), bindings[1].(float64), "the 2nd binding should be 64.56")
		assert.Equal(t, int(10), bindings[2].(int), "the 3rd binding should be 10")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 1, len(rows), "the return value should has 1 row")
	if len(rows) == 1 {
		assert.Equal(t, "john@yao.run", rows[0]["email"].(string), "the email of first row should be john@yao.run")
		assert.Equal(t, int64(1), rows[0]["id"].(int64), "the email of first row should be 1")
		assert.Equal(t, "WAITING", rows[0]["status"].(string), "the email of first row should be WAITING")
		if unit.DriverIs("sqlite3") {
			assert.Equal(t, "2021-03-25 00:21:16", rows[0]["created_at"].(string), "the email of first row should be WAITING")
		} else {
			assert.Equal(t, "2021-03-25T00:21:16", rows[0]["created_at"].(time.Time).Format("2006-01-02T15:04:05"), "the email of first row should be WAITING")
		}

	}
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

	// select * from `table_test_where` where `email` like ? and (`vote` > ? and `name` = ? and (`created_at` > ? and `created_at` < ?)) and `score` > ?
	// select * from "table_test_where" where "email" like $1 and ("vote" > $2 and "name" = $3 and ("created_at" > $4 and "created_at" < $5)) and "score" > $6

	// checking sql
	sql := qb.ToSQL()
	if unit.DriverIs("postgres") {
		assert.Equal(t, `select * from "table_test_where" where "email" like $1 and ("vote" > $2 and "name" = $3 and ("created_at" > $4 and "created_at" < $5)) and "score" > $6`, sql, "the query sql not equal")
	} else {
		assert.Equal(t, "select * from `table_test_where` where `email` like ? and (`vote` > ? and `name` = ? and (`created_at` > ? and `created_at` < ?)) and `score` > ?", sql, "the query sql not equal")
	}

	// checking bindings
	bindings := qb.GetBindings()
	assert.Equal(t, 6, len(bindings), "the bindings should have 3 items")
	if len(bindings) == 6 {
		assert.Equal(t, "%@yao.run", bindings[0].(string), "the 1st binding should be %@yao.run")
		assert.Equal(t, int(10), bindings[1].(int), "the 2nd binding should be 10")
		assert.Equal(t, "Ken", bindings[2].(string), "the 3rd binding should be Ken")
		assert.Equal(t, "2021-03-25 08:00:00", bindings[3].(string), "the 4th binding should be 2021-03-25 08:00:00")
		assert.Equal(t, "2021-03-25 19:00:00", bindings[4].(string), "the 5th binding should be 2021-03-25 19:00:00")
		assert.Equal(t, int(5), bindings[5].(int), "the 5th binding should be 2021-03-25 19:00:00")
	}

	// checking result
	rows := qb.MustGet()
	assert.Equal(t, 1, len(rows), "the return value should has 1 row")
	if len(rows) == 1 {
		assert.Equal(t, "ken@yao.run", rows[0]["email"].(string), "the email of first row should be ken@yao.run")
	}
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
