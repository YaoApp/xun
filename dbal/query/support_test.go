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

func TestSupportIsBoolean(t *testing.T) {
	builder := getTestBuilderInstance()
	assert.True(t, builder.isBoolean("and"), "The return value should be true")
	assert.True(t, builder.isBoolean("or"), "The return value should be true")
	assert.False(t, builder.isBoolean("not"), "The return value should be false")
}

func TestFilterNilBindingsRemovesNil(t *testing.T) {
	values := []interface{}{1, nil, "hello", nil, 3}
	result := filterNilBindings(values)
	assert.Equal(t, []interface{}{1, "hello", 3}, result)
}

func TestFilterNilBindingsTypedNil(t *testing.T) {
	var p *string
	values := []interface{}{"a", p, "b"}
	result := filterNilBindings(values)
	assert.Equal(t, []interface{}{"a", "b"}, result)
}

func TestFilterNilBindingsAllNil(t *testing.T) {
	values := []interface{}{nil, nil, nil}
	result := filterNilBindings(values)
	assert.Empty(t, result)
}

func TestFilterNilBindingsNoNil(t *testing.T) {
	values := []interface{}{1, "two", 3.0}
	result := filterNilBindings(values)
	assert.Equal(t, []interface{}{1, "two", 3.0}, result)
}

func TestFilterNilBindingsEmpty(t *testing.T) {
	result := filterNilBindings([]interface{}{})
	assert.Empty(t, result)
}
