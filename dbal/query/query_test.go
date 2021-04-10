package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestQueryMustGet(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	rows := qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustGet()

	assert.Equal(t, 4, len(rows), "the return rows should have 4 items")
	if len(rows) == 4 {
		assert.Equal(t, "96.32", fmt.Sprintf("%.2f", rows[0].Get("score")), "the return value should be true")
		assert.Equal(t, "64.56", fmt.Sprintf("%.2f", rows[1].Get("score")), "the return value should be true")
		assert.Equal(t, "99.27", fmt.Sprintf("%.2f", rows[2].Get("score")), "the return value should be true")
		assert.Equal(t, "48.12", fmt.Sprintf("%.2f", rows[3].Get("score")), "the return value should be true")
	}
}

func TestQueryMustGetBind(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	type Item struct {
		ID            int64
		Email         string
		Score         float64
		Vote          int
		ScoreGrade    xun.N
		CreatedAt     xun.T
		UpdatedAt     xun.T
		PaymentStatus string `json:"status"`
		Extra         string
	}

	rows := []Item{}
	qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "score", "vote", "status", "score_grade", "created_at", "updated_at").
		OrderBy("id").
		MustGet(&rows)

	assert.Equal(t, 4, len(rows), "the return rows should have 4 items")
	if len(rows) == 4 {
		assert.Equal(t, "96.32", fmt.Sprintf("%.2f", rows[0].Score), "the return value should be true")
		assert.Equal(t, 10, rows[0].Vote, "the return value should be true")
		assert.Equal(t, 99.27, rows[0].ScoreGrade.MustToFixed(2), "the return value should be true")
		assert.Equal(t, "WAITING", rows[0].PaymentStatus, "the return value should be true")
		assert.Equal(t, 2021, rows[0].CreatedAt.MustToTime().Year(), "the return value should be true")
		assert.True(t, rows[0].UpdatedAt.IsNull(), "the return value should be true")

		assert.Equal(t, "64.56", fmt.Sprintf("%.2f", rows[1].Score), "the return value should be true")
		assert.Equal(t, 5, rows[1].Vote, "the return value should be true")
		assert.Equal(t, "99.27", fmt.Sprintf("%.2f", rows[2].Score), "the return value should be true")
		assert.Equal(t, 125, rows[2].Vote, "the return value should be true")
		assert.Equal(t, "48.12", fmt.Sprintf("%.2f", rows[3].Score), "the return value should be true")
		assert.Equal(t, 6, rows[3].Vote, "the return value should be true")
	}
}

func TestQueryMustExistsTrue(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	res := qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustExists()

	assert.True(t, res, "the return value should be true")
}

func TestQueryMustExistsFalse(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	res := qb.From("table_test_query as t").
		Where("email", "like", "%@iqka.com").
		OrderBy("id").
		MustExists()

	assert.False(t, res, "the return value should be false")
}

func TestQueryMustDoesntExistTrue(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	res := qb.From("table_test_query as t").
		Where("email", "like", "%@iqka.com").
		OrderBy("id").
		MustDoesntExist()

	assert.True(t, res, "the return value should be true")
}

func TestQueryMustDoesntExistFalse(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	res := qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustDoesntExist()

	assert.False(t, res, "the return value should be false")
}

func TestQueryMustFirst(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	row := qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		OrderBy("id").
		MustFirst()

	assert.Equal(t, "john@yao.run", row.Get("email"), "the email should be true")
}

func TestQueryMustFirstBind(t *testing.T) {
	type Item struct {
		ID            int64
		Email         string
		Score         float64
		Vote          int
		ScoreGrade    xun.N
		CreatedAt     xun.T
		UpdatedAt     xun.T
		PaymentStatus string `json:"status"`
		Extra         string
	}

	NewTableForQueryTest()
	qb := getTestBuilder()

	row := Item{}
	qb.From("table_test_query as t").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "score", "vote", "status", "score_grade", "created_at", "updated_at").
		OrderBy("id").
		MustFirst(&row)

	assert.Equal(t, "john@yao.run", row.Email, "the email should be true")
}

func TestQueryMustFirstEmpty(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	row := qb.From("table_test_query as t").
		Where("email", "like", "%@yaorun").
		OrderBy("id").
		MustFirst()

	assert.True(t, row.IsEmpty(), "the return row should be empty")
}

func TestQueryMustFirstError(t *testing.T) {
	NewTableForQueryTest()
	qb := getTestBuilder()
	assert.Panics(t, func() {
		qb.From("table_test_query as t").
			Where("ping", "like", "%@yaorun").
			OrderBy("id").
			MustFirst()
	})
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
