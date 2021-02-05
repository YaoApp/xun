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
	builder := getTestBuilder()
	builder.CreateT("table_test_builder", func(blueprint *Blueprint) {

	})
}
