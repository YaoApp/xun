package schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func getTestBuilder() *schema.Builder {
	defer unit.Catch()
	name := os.Getenv("XUN_UNIT_DSN")
	dsn := unit.DSN(name)
	return schema.NewBuilderByDSN(name, dsn)
}

func TestCreate(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.Create("table_test_builder", func(table *schema.Blueprint) {
		table.String("name", 20).Unique()
		table.String("unionid", 128).Unique()
	})
	assert.True(t, true, "the table should be true")
}

func TestExistsTrue(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table := builder.Table("table_test_builder")
	assert.True(t, table.Exists(), "should return true")
}

func TestDrop(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.Drop("table_test_builder")
}

func TestExistsFalse(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table := builder.Table("table_test_builder")
	assert.True(t, !table.Exists(), "should return false")
}
