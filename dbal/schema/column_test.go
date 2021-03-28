package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

func TestColumnGetColumn(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	NewColumnsForTest(table)

	col1 := table.GetColumn("field1")
	assert.True(t, col1 != nil, "the column field1 should not be nil")
	if col1 != nil {
		assert.Equal(t, "field1", col1.Name, "the column name should be field1")
		assert.Equal(t, 20, utils.IntVal(col1.Length), "the column field1 length should be 20")
	}

	col2 := table.GetColumn("field2")
	assert.True(t, col2 != nil, "the column field2 should not be nil")
	if col2 != nil {
		assert.Equal(t, "field2", col2.Name, "the column name should be field2")
		assert.Equal(t, 40, utils.IntVal(col2.Length), "the column field2 length should be 20")
	}
}

func TestColumnHasColumn(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	NewColumnsForTest(table)
	assert.True(t, table.HasColumn("field1"), "the table test should have the field1 column")
	assert.True(t, table.HasColumn("field2"), "the table test should have the field2 column")
	assert.False(t, table.HasColumn("field3"), "the table test should not have the field3 column")
}

func TestColumnRenameColumn(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	NewTableForColumnTest()
	builder.MustAlterTable("table_test_column", func(table Blueprint) {
		table.RenameColumn("field2", "re_field2")
	})
	table := builder.MustGetTable("table_test_column")
	assert.True(t, table.HasColumn("field1"), "the table table_test_column should have the field1 column")
	assert.True(t, table.HasColumn("re_field2"), "the table table_test_column should have the re_field2 column")
	assert.False(t, table.HasColumn("field2"), "the table table_test_column should not have the field2 column")
}

func TestColumnRenameColumnFail(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	NewTableForColumnTest()
	assert.Panics(t, func() {
		builder.MustAlterTable("table_test_column", func(table Blueprint) {
			table.RenameColumn("field2", "field1")
		})
	})
}

func TestColumnDropColumn(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	defer unit.Catch()
	builder := getTestBuilder()
	NewTableForColumnTest()
	builder.MustAlterTable("table_test_column", func(table Blueprint) {
		table.DropColumn("field2")
	})
	table := builder.MustGetTable("table_test_column")
	assert.True(t, table.HasColumn("field1"), "the table table_test_column should have the field1 column")
	assert.False(t, table.HasColumn("field2"), "the table table_test_column should not have the field2 column")
}

func TestColumnSetLength(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	var col *Column
	col = table.String("col", 20)
	assert.Equal(t, 20, utils.IntVal(col.Length), "the column length should be 20")

	col = table.String("col", 65536)
	assert.Equal(t, 200, utils.IntVal(col.Length), "the column length should be 200")

	col = table.String("col", 0)
	assert.Equal(t, 200, utils.IntVal(col.Length), "the column length should be 200")

	col = table.String("col", 65535)
	assert.Equal(t, 65535, utils.IntVal(col.Length), "the column length should be 65535")

	col = table.String("col", 65535).SetLength(256)
	assert.Equal(t, 256, utils.IntVal(col.Length), "the column length should be 256")

	col = table.BigIncrements("col").SetLength(200)
	assert.True(t, col.Length == nil, "the column length should be nil")
}

func TestColumnSetPrecision(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	var col *Column
	col = table.Float("col", 10, 2)
	assert.Equal(t, 10, utils.IntVal(col.Precision), "the column precision should be 10")
	assert.Equal(t, 2, utils.IntVal(col.Scale), "the column scale should be 2")

	col = table.Float("col", 23, 2)
	assert.Equal(t, 23, utils.IntVal(col.Precision), "the column precision should be 23")
	assert.Equal(t, 0, utils.IntVal(col.Scale), "the column scale should be 0")

	col = table.Float("col", 24, 2)
	assert.Equal(t, 10, utils.IntVal(col.Precision), "the column precision should be 10")
	assert.Equal(t, 2, utils.IntVal(col.Scale), "the column scale should be 2")

	col = table.Float("col", 0, 2)
	assert.Equal(t, 10, utils.IntVal(col.Precision), "the column precision should be 10")
	assert.Equal(t, 2, utils.IntVal(col.Scale), "the column scale should be 2")

	col = table.Float("col", 0, 23)
	assert.Equal(t, 10, utils.IntVal(col.Precision), "the column precision should be 10")
	assert.Equal(t, 2, utils.IntVal(col.Scale), "the column scale should be 2")

	col = table.Float("col", 10, 4).SetPrecision(12)
	assert.Equal(t, 12, utils.IntVal(col.Precision), "the column precision should be 12")
	assert.Equal(t, 4, utils.IntVal(col.Scale), "the column scale should be 4")

	col = table.Float("col").SetScale(6).SetPrecision(20)
	assert.Equal(t, 17, utils.IntVal(col.Precision), "the column precision should be 12")
	assert.Equal(t, 6, utils.IntVal(col.Scale), "the column scale should be 4")

	col = table.BigIncrements("col").SetPrecision(200)
	assert.True(t, col.Precision == nil, "the column precision should be nil")

}

func TestColumnSetScale(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	var col *Column

	col = table.Float("col", 10, 4).SetScale(6)
	assert.Equal(t, 10, utils.IntVal(col.Precision), "the column precision should be 10")
	assert.Equal(t, 6, utils.IntVal(col.Scale), "the column scale should be 4")

	col = table.BigIncrements("col").SetScale(200)
	assert.True(t, col.Scale == nil, "the column scale should be nil")
}

func TestColumnSetDateTimePrecision(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	var col *Column
	col = table.DateTime("col").SetDateTimePrecision(6)
	assert.Equal(t, 6, utils.IntVal(col.DateTimePrecision), "the column DateTimePrecision should be 6")

	col = table.DateTime("col").SetDateTimePrecision(12)
	assert.Equal(t, 0, utils.IntVal(col.DateTimePrecision), "the column DateTimePrecision should be 0")

	col = table.DateTime("col").SetDateTimePrecision(5)
	assert.Equal(t, 5, utils.IntVal(col.DateTimePrecision), "the column DateTimePrecision should be 5")

	col = table.Date("col").SetDateTimePrecision(5)
	assert.True(t, col.DateTimePrecision == nil, "the column DateTimePrecision should be nil")
}

func TestColumnSetComment(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	var col *Column
	col = table.String("col").SetComment("This is a col")
	assert.Equal(t, "This is a col", utils.StringVal(col.Comment), "the column Comment should be  \"This is a col\"")
	col = table.String("col")
	assert.True(t, col.Comment == nil, "the column Comment should be nil")
}

func TestColumnHasIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	NewTableForColumnTest()
	table := builder.MustGetTable("table_test_column")
	col := table.GetColumn("field1")
	assert.True(t, col != nil, "the column field1 should be exists")
	res := col.HasIndex("field1_index")
	assert.True(t, res, "the column field1 should have the index field1_index")

	// @todo:  postgres does not have the index field1_field2. it should be fixed at next version.
	// res = col.HasIndex("field1_field2")
}

func TestColumnUniqueAndIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_column")
	builder.MustCreateTable("table_test_column", func(table Blueprint) {
		table.ID("id")
		table.String("field1").Index().Index()
		table.String("field2").Unique().Unique()
		table.AddIndex("field1_field2", "field1", "field2")
	})

	assert.True(t, builder.MustHasTable("table_test_column"), "the table table_test_column should be created")
	if builder.MustHasTable("table_test_column") {
		table := builder.MustGetTable("table_test_column")
		col1 := table.GetColumn("field1")
		col2 := table.GetColumn("field2")
		assert.True(t, col1.HasIndex("field1_index"), "the column field1 should have the index field1_index")
		assert.True(t, col2.HasIndex("field2_unique"), "the column field1 should have the index field2_unique")
	}
}

// clean the test data
func TestColumnClean(t *testing.T) {
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_column")
}

func NewColumnsForTest(table *Table) {
	col1 := table.String("field1", 20)
	col2 := table.String("field2", 40)
	table.pushColumn(col1)
	table.pushColumn(col2)
}

func NewTableForColumnTest() {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_column")
	builder.MustCreateTable("table_test_column", func(table Blueprint) {
		table.ID("id")
		table.String("field1").Index()
		table.String("field2")
		table.String("field3").SetDefault("DefaultValue3")
		table.AddIndex("field1_field2", "field1", "field2")
	})
}
