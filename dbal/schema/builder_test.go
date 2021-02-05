package schema

import (
	"os"
	"testing"
)

func getTestBuilder() *Builder {
	dsn := os.Getenv("XUN_UNIT_MYSQL_DSN")
	return NewBuilderByDSN("mysql", dsn)
}

func TestCreate(t *testing.T) {
	schema := getTestBuilder()
	schema.Create("table_test_builder", func(table *Blueprint) {
		table.String("name", 20).Unique()
		table.String("unionid", 128).Unique()
	})
}
