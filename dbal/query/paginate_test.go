package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestPaginateBasic(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	paginateor := qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email", "vote", "score", "status").
		OrderBy("vote", "desc").
		OrderBy("score").
		Limit(20).Offset(100).
		MustPaginate(2, 1)

	// cheking paginateor
	assert.Equal(t, 4, paginateor.Total, "The records total should be 4")
	assert.Equal(t, 2, paginateor.PageSize, "The page size should be 2")
	assert.Equal(t, 2, paginateor.TotalPages, "The total pages should be 2")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, 2, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// cheking items
	if len(paginateor.Items) == 2 {
		assert.Equal(t, int64(3), paginateor.Items[0].(xun.R).Get("id"), "The first row id should be 3")
		assert.Equal(t, int64(1), paginateor.Items[1].(xun.R).Get("id"), "The second row id should be 1")
	}
	// cheking next page
	paginateor = qb.MustPaginate(2, 2)
	assert.Equal(t, 4, paginateor.Total, "The records total should be 4")
	assert.Equal(t, 2, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, -1, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, 1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// cheking items
	if len(paginateor.Items) == 2 {
		assert.Equal(t, int64(4), paginateor.Items[0].(xun.R).Get("id"), "The first row id should be 4")
		assert.Equal(t, int64(2), paginateor.Items[1].(xun.R).Get("id"), "The second row id should be 2")
	}
}

func TestPaginateBasicBind(t *testing.T) {
	NewTableForPaginateTest()
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

	items := []Item{}
	paginateor := qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "vote", "score", "status", "created_at").
		OrderBy("vote", "desc").
		OrderBy("score").
		Limit(20).Offset(100).
		MustPaginate(2, 1, &items)

	// cheking paginateor
	assert.Equal(t, 4, paginateor.Total, "The records total should be 4")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, 2, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// cheking items
	if len(paginateor.Items) == 2 {
		assert.Equal(t, int64(3), paginateor.Items[0].(Item).ID, "The first row id should be 3")
		assert.Equal(t, int64(1), paginateor.Items[1].(Item).ID, "The second row id should be 1")
	}

	// cheking next page
	items = []Item{}
	paginateor = qb.MustPaginate(2, 2, &items)
	assert.Equal(t, 4, paginateor.Total, "The records total should be 4")
	assert.Equal(t, 2, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, -1, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, 1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// cheking items
	if len(paginateor.Items) == 2 {
		assert.Equal(t, int64(4), paginateor.Items[0].(Item).ID, "The first row id should be 4")
		assert.Equal(t, int64(2), paginateor.Items[1].(Item).ID, "The second row id should be 2")
	}
}

func TestPaginateUnion(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "status").
		OrderBy("id").
		UnionAll(func(qb Query) {
			qb.Table("table_test_paginate_t2").
				Select("id", "name", "status")
		})
	paginateor := qb.MustPaginate(2, 1)

	// cheking paginateor
	assert.Equal(t, 8, paginateor.Total, "The records total should be 8")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, 2, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// cheking items
	if len(paginateor.Items) == 2 {
		assert.Equal(t, int64(1), paginateor.Items[0].(xun.R).Get("id"), "The first row id should be 3")
		assert.Equal(t, int64(2), paginateor.Items[1].(xun.R).Get("id"), "The second row id should be 1")
	}

	// @Todo: union query has bug. not fixed.
	// ******
	// cheking next page
	// paginateor = qb.MustPaginate(2, 2)
	// utils.Println(paginateor)

	// assert.Equal(t, 8, paginateor.Total, "The records total should be 4")
	// assert.Equal(t, 2, paginateor.CurrentPage, "The current page should be 1")
	// assert.Equal(t, 3, paginateor.NextPage, "The next page should be 2")
	// assert.Equal(t, 1, paginateor.PreviousPage, "The previous page should be -1")
	// assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

	// // cheking items
	// if len(paginateor.Items) == 2 {
	// 	assert.Equal(t, int64(4), paginateor.Items[0].(xun.R).Get("id"), "The first row id should be 4")
	// 	assert.Equal(t, int64(2), paginateor.Items[1].(xun.R).Get("id"), "The second row id should be 2")
	// }
}

func TestPaginateWithGroup(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("status").
		GroupBy("status").
		Limit(20).Offset(100)

	paginateor := qb.MustPaginate(2, 1)

	// cheking paginateor
	assert.Equal(t, 3, paginateor.Total, "The records total should be 8")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, 2, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")

}

func TestPaginateWithGroupEmpty(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@hello.run").
		Select("status").
		GroupBy("status").
		Limit(20).Offset(100)

	paginateor := qb.MustPaginate(2, 1)

	// cheking paginateor
	assert.Equal(t, 0, paginateor.Total, "The records total should be 8")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, -1, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 0, len(paginateor.Items), "The items count should be 2")
}

func TestPaginateWithGroupAndJoin(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate as t1").
		CrossJoin("table_test_paginate_t2 as t2").
		Select("t1.status").
		Where("email", "like", "%@yao.run").
		GroupBy("t1.status").
		Limit(20).Offset(100)

	paginateor := qb.MustPaginate(2, 1)

	// cheking paginateor
	assert.Equal(t, 3, paginateor.Total, "The records total should be 8")
	assert.Equal(t, 1, paginateor.CurrentPage, "The current page should be 1")
	assert.Equal(t, 2, paginateor.NextPage, "The next page should be 2")
	assert.Equal(t, -1, paginateor.PreviousPage, "The previous page should be -1")
	assert.Equal(t, 2, len(paginateor.Items), "The items count should be 2")
}

func TestPaginateChunk(t *testing.T) {
	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "name", "email", "vote", "score", "status").
		OrderByDesc("id")
	hits := 0
	IDs := []int64{}
	total := qb.MustCount()
	assert.Equal(t, 4, int(total), "The query should have 4 items")
	qb.MustChunk(2, func(items []interface{}, page int) error {
		hits = hits + len(items)
		for _, item := range items {
			IDs = append(IDs, item.(xun.R).Get("id").(int64))
		}
		return nil
	})
	assert.Equal(t, hits, int(total), "The chunk items hits should be equal total")
	assert.Equal(t, []int64{4, 3, 2, 1}, IDs, "The chunk id of items ids should be []int64{4,3,2,1}")
}

func TestPaginateChunkWithBind(t *testing.T) {

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

	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "vote", "score", "status").
		OrderByDesc("id")
	hits := 0
	IDs := []int64{}
	total := qb.MustCount()
	assert.Equal(t, 4, int(total), "The query should have 4 items")
	qb.MustChunk(2, func(items []interface{}, page int) error {
		hits = hits + len(items)
		for _, item := range items {
			IDs = append(IDs, item.(Item).ID)
		}
		return nil
	}, &[]Item{})
	assert.Equal(t, hits, int(total), "The chunk items hits should be equal total")
	assert.Equal(t, []int64{4, 3, 2, 1}, IDs, "The chunk id of items ids should be []int64{4,3,2,1}")

}

func TestPaginateChunkWithBindError(t *testing.T) {

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

	NewTableForPaginateTest()
	qb := getTestBuilder()
	qb.Table("table_test_paginate").
		Where("email", "like", "%@yao.run").
		Select("id", "email", "vote", "score", "status").
		OrderByDesc("id")

	assert.PanicsWithError(t, "The given binding var shoule be a slice pointer", func() {
		qb.MustChunk(2, func(items []interface{}, page int) error {
			return nil
		}, &Item{})
	})
}

// clean the test data
func TestPaginateClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_paginate")
	builder.DropTableIfExists("table_test_paginate_t2")
}

func NewTableForPaginateTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_paginate")
	builder.MustCreateTable("table_test_paginate", func(table schema.Blueprint) {
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
	qb.Table("table_test_paginate").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})

	builder.DropTableIfExists("table_test_paginate_t2")
	builder.MustCreateTable("table_test_paginate_t2", func(table schema.Blueprint) {
		table.ID("id")
		table.BigInteger("t1_id")
		table.String("name", 300)
		table.Enum("status", []string{"WAITING", "PENDING", "DONE"}).SetDefault("WAITING")
		table.Timestamps()
		table.SoftDeletes()
	})
	qb.Table("table_test_paginate_t2").Insert([]xun.R{
		{"t1_id": 1, "name": "Emma", "status": "WAITING", "created_at": "2021-03-27 00:00:16"},
		{"t1_id": 2, "name": "Ava", "status": "PENDING", "created_at": "2021-03-27 08:13:23"},
		{"t1_id": 3, "name": "Amelia", "status": "DONE", "created_at": "2021-03-27 09:46:21"},
		{"t1_id": 4, "name": "Elizabeth", "status": "DONE", "created_at": "2021-03-27 14:00:22"},
	})
}
