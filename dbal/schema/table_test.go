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
