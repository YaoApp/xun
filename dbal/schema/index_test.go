package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestIndexGetIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	table.pushIndex(NewIndexForTest(table, "index_1"))
	index := table.GetIndex("index_1")
	assert.Equal(t, 2, len(index.Columns), "the index should have 2 columns")
	assert.Equal(t, "index", index.Type, "the index type shoule be index")
	assert.Equal(t, "index_1", index.Name, "the index name should be index_1")
	assert.Equal(t, "test", table.Name, "the table name should be test")
}

func TestIndexHasIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	index := table.newIndex("index_1")
	table.pushIndex(index)
	assert.True(t, table.HasIndex("index_1"), "the table should have the index_1 index")
	assert.False(t, table.HasIndex("index_2"), "the table should not have the index_2 index")
}

func TestIndexAddIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.MustDropTableIfExists("table_test_index")
	builder.MustCreateTable("table_test_index", func(table Blueprint) {
		table.ID("id")
		table.String("field1")
		table.String("field2")
		table.AddIndex("field1_field2", "field1", "field2")
	})

	table := builder.MustGetTable("table_test_index")
	assert.True(t, table.HasIndex("field1_field2"), "the table should have the field1_field2 index")
	if table.HasIndex("field1_field2") {
		index := table.GetIndex("field1_field2")
		assert.Equal(t, "index", index.Type, "the type of field1_field2 index should be 'index'")
	}

	// AddIndex Fail
	err := builder.AlterTable("table_test_index", func(table Blueprint) {
		table.AddIndex("field1_field2", "id", "field1", "field2")
	})
	assert.False(t, err == nil, "The return error should not be nil")
}

func TestIndexAddUnique(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_index")
	builder.MustCreateTable("table_test_index", func(table Blueprint) {
		table.ID("id")
		table.String("field1", 40)
		table.String("field2", 40)
		table.AddUnique("field1_field2", "field1", "field2")
	})
	table := builder.MustGetTable("table_test_index")
	assert.True(t, table.HasIndex("field1_field2"), "the table should have the field1_field2 index")
	if table.HasIndex("field1_field2") {
		index := table.GetIndex("field1_field2")
		assert.Equal(t, "unique", index.Type, "the type of field1_field2 index should be 'index'")
	}

	// AddUnique Fail
	err := builder.AlterTable("table_test_index", func(table Blueprint) {
		table.AddUnique("field1_field2", "id", "field1", "field2")
	})
	assert.False(t, err == nil, "The return error should not be nil")
}

func TestIndexRenameIndex(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	defer unit.Catch()
	builder := getTestBuilder()
	TestIndexAddIndex(t)
	builder.MustAlterTable("table_test_index", func(table Blueprint) {
		table.RenameIndex("field1_field2", "re_field1_field2")
	})
	table := builder.MustGetTable("table_test_index")
	assert.True(t, table.HasIndex("re_field1_field2"), "the table should have the re_field1_field2 index")
	assert.False(t, table.HasIndex("field1_field2"), "the table should have not the field1_field2 index")

	// RenameIndex Fail
	err := builder.AlterTable("table_test_index", func(table Blueprint) {
		table.AddIndex("field1_field2", "field2")
		table.RenameIndex("re_field1_field2", "field1_field2")
	})
	assert.False(t, err == nil, "The return error should not be nil")
}

func TestIndexDropIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	TestIndexAddIndex(t)
	builder.MustAlterTable("table_test_index", func(table Blueprint) {
		table.DropIndex("field1_field2")
	})
	table := builder.MustGetTable("table_test_index")
	assert.False(t, table.HasIndex("field1_field2"), "the table should have not the field1_field2 index")
}

// clean the test data
func TestIndexClean(t *testing.T) {
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_index")
}

func NewIndexForTest(table *Table, name string) *Index {
	col1 := table.String("field1", 20)
	col2 := table.String("field2", 20)
	return table.newIndex(name, col1, col2)
}
