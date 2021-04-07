package query

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestDebugDD(t *testing.T) {
	if os.Getenv("BE_EXIT") == "1" {
		NewTableForWhereTest()
		qb := getTestBuilder()
		qb.Table("table_test_where").Where("id", 1)
		qb.DD()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDebugDD")
	cmd.Env = append(os.Environ(), "BE_EXIT=1")
	bytes, err := cmd.Output()
	out := string(bytes)
	assert.Nil(t, err, "the command should be executed success")
	assert.True(t, strings.Contains(out, "select"), "the command return value should be have  select...")
	assert.True(t, strings.Contains(out, `"id": 1`), "the command return value should be have  ...id")
}

func TestDebugDDFail(t *testing.T) {
	if os.Getenv("BE_EXIT") == "1" {
		NewTableForWhereTest()
		qb := getTestBuilder()
		qb.Table("table_test_where").WhereRaw(`something error`)
		qb.DD()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDebugDDFail")
	cmd.Env = append(os.Environ(), "BE_EXIT=1")
	bytes, err := cmd.Output()
	out := string(bytes)
	assert.Nil(t, err, "the command should be executed success")
	assert.True(t, strings.Contains(out, "select"), "the command return value should be have select...")
	assert.True(t, strings.Contains(out, `[]`), "the command return value should be have []...")
	assert.True(t, strings.Contains(out, `runtime/debug.Stack`), "the command return value should be have runtime/debug.Stack...")
}

// clean the test data
func TestDebugClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_debug")
}

func NewTableForDebugTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_debug")
	builder.MustCreateTable("table_test_debug", func(table schema.Blueprint) {
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
	qb.Table("table_test_debug").Insert([]xun.R{
		{"email": "john@yao.run", "name": "John", "vote": 10, "score": 96.32, "score_grade": 99.27, "status": "WAITING", "created_at": "2021-03-25 00:21:16"},
		{"email": "lee@yao.run", "name": "Lee", "vote": 5, "score": 64.56, "score_grade": 99.27, "status": "PENDING", "created_at": "2021-03-25 08:30:15"},
		{"email": "ken@yao.run", "name": "Ken", "vote": 125, "score": 99.27, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 09:40:23"},
		{"email": "ben@yao.run", "name": "Ben", "vote": 6, "score": 48.12, "score_grade": 99.27, "status": "DONE", "created_at": "2021-03-25 18:15:29"},
	})
}
