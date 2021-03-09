package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

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
