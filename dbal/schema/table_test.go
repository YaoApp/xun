package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

func TestTableGetPrefix(t *testing.T) {
	builder := newBuilder(unit.Driver(), unit.DSN())
	builder.SetOption(&dbal.Option{
		Prefix: "xun_",
	})

	builder.DropTableIfExists("table_test_table")
	builder.MustCreateTable("table_test_table", func(table Blueprint) {
		table.ID("id")
		table.String("name", 200)
	})

	tables := builder.MustGetTables()
	assert.True(t, utils.StringHave(tables, "table_test_table"), "the talbe prefix should be xun_")

	table := builder.MustGetTable("table_test_table")
	assert.Equal(t, "xun_", table.GetPrefix(), "the talbe prefix should be xun_")
}

func TestTableCreateTemporary(t *testing.T) {
	option := dbal.CreateTableOption{
		Temporary: true,
		Engine:    "MEMORY",
	}

	// Use the memory engine
	builder := newBuilder(unit.Driver(), unit.DSN())
	builder.DropTableIfExists("table_test_table_temp")
	builder.MustCreateTable("table_test_table_temp", func(table Blueprint) {
		table.ID("id")
		table.String("name", 200)
	}, option)

	has, err := builder.HasTable("table_test_table_temp")
	assert.Nil(t, err, "the table should be created")

	// Get the table on mysql or sqlite3, PostgreSQL not supported, currently
	if unit.DriverIs("mysql") || unit.DriverIs("sqlite3") {
		table, err := builder.GetTable("table_test_table_temp")
		assert.Nil(t, err, "the table should be created")
		assert.Equal(t, "table_test_table_temp", table.GetName(), "the table name should be table_test_table_temp")
	}

	// Driver is mysql or sqlite3, PostgreSQL will return false
	// fully supported in the future
	if unit.DriverIs("mysql") || unit.DriverIs("sqlite3") {
		assert.True(t, has, "the table should be created")
	}

	// Use temporary table
	option = dbal.CreateTableOption{Temporary: true}
	builder.DropTableIfExists("table_test_table_temp")
	builder.MustCreateTable("table_test_table_temp", func(table Blueprint) {
		table.ID("id")
		table.String("name", 200)
	}, option)

	has, err = builder.HasTable("table_test_table_temp")
	assert.Nil(t, err, "the table should be created")

	if unit.DriverIs("mysql") {
		// HasTable can not check temporary table, driver mysql will return false
		// It will be fixed in the future
		assert.False(t, has, "the table should be created")
	}
}
