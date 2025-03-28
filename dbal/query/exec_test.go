package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestExec(t *testing.T) {
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

	// Exec
	if unit.Is("postgres") {
		res, err := qb.Exec("update table_test_query set score = 100 where email like $1", "%@yao.run")
		assert.Nil(t, err, "the error should be nil")
		affected, err := res.RowsAffected()
		assert.Nil(t, err, "the error should be nil")
		assert.Equal(t, int64(4), affected, "the rows affected should be 4")
	} else {
		res, err := qb.Exec("update table_test_query set score = 100 where email like ?", "%@yao.run")
		assert.Nil(t, err, "the error should be nil")
		affected, err := res.RowsAffected()
		assert.Nil(t, err, "the error should be nil")
		assert.Equal(t, int64(4), affected, "the rows affected should be 4")
	}
}
func TestExecWrite(t *testing.T) {
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

	// Exec
	if unit.Is("postgres") {
		res, err := qb.ExecWrite("update table_test_query set score = 100 where email like $1", "%@yao.run")
		assert.Nil(t, err, "the error should be nil")
		affected, err := res.RowsAffected()
		assert.Nil(t, err, "the error should be nil")
		assert.Equal(t, int64(4), affected, "the rows affected should be 4")
	} else {
		res, err := qb.ExecWrite("update table_test_query set score = 100 where email like ?", "%@yao.run")
		assert.Nil(t, err, "the error should be nil")
		affected, err := res.RowsAffected()
		assert.Nil(t, err, "the error should be nil")
		assert.Equal(t, int64(4), affected, "the rows affected should be 4")
	}
}
