package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestIndexNewIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	col1 := table.String("field1", 20)
	col2 := table.String("field2", 20)
	index := table.NewIndex("index_1", col1, col2)
	assert.Equal(t, 2, len(index.Columns), "the index should have 2 columns")
	assert.Equal(t, "index", index.Type, "the index type shoule be index")
	assert.Equal(t, "index_1", index.Name, "the index name should be index_1")
	assert.Equal(t, "test", table.Name, "the table name should be test")
}

func TestIndexPushIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	index := NewIndexForTest(table, "index_1")
	table.PushIndex(index)
	assert.True(t, table.HasIndex("index_1"), "the table should has the index_1 index")
}

func TestIndexGetIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	table.PushIndex(NewIndexForTest(table, "index_1"))
	index := table.GetIndex("index_1")
	assert.Equal(t, 2, len(index.Columns), "the index should have 2 columns")
	assert.Equal(t, "index", index.Type, "the index type shoule be index")
	assert.Equal(t, "index_1", index.Name, "the index name should be index_1")
	assert.Equal(t, "test", table.Name, "the table name should be test")
}

func TestIndexIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	index := table.Index("index_1")
	assert.Equal(t, "index_1", index.Name, "the index name should be index_1")
	assert.Equal(t, "test", table.Name, "the table name should be test")
}

func TestIndexHasIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	index := table.Index("index_1")
	table.PushIndex(index)
	assert.True(t, table.HasIndex("index_1"), "the table should have the index_1 index")
	assert.False(t, table.HasIndex("index_2"), "the table should not have the index_2 index")
}

func TestIndexPutIndex(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	table := NewTable("test", builder)
	index := table.Index("index_1")
	table.PushIndex(index)
	assert.True(t, table.HasIndex("index_1"), "the table should have the index_1 index")
	assert.False(t, table.HasIndex("index_2"), "the table should not have the index_2 index")
}

func NewIndexForTest(table *Table, name string) *Index {
	col1 := table.String("field1", 20)
	col2 := table.String("field2", 20)
	return table.NewIndex(name, col1, col2)
}
