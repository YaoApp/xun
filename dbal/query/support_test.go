package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportIsQueryable(t *testing.T) {
	NewTableForWhereTest()
	qb := getTestBuilder()
	builder := getTestBuilderInstance()
	var empty interface{} = nil
	assert.False(t, builder.isQueryable(empty), "The nil interface{} should not be queryable")
	assert.False(t, builder.isQueryable(func() {}), "The func(){} should not be queryable")
	assert.True(t, builder.isQueryable(func(qb Query) {}), "The func(qb Query) {} should be queryable")
	assert.True(t, builder.isQueryable(builder), "The builder instance should be queryable")
	assert.True(t, builder.isQueryable(qb), "The Query interface should be queryable")
}
