package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
)

func newTestSQL() SQL {
	return NewSQL(&Quoter{})
}

func TestWhereJsoncontainsMySQL(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsMySQLNot(t *testing.T) {
	g := newTestSQL()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "not JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsMySQLOffset(t *testing.T) {
	g := newTestSQL()
	offset := 3
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"test"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "JSON_CONTAINS(`tags`, ?)", result)
	assert.Equal(t, 4, offset)
}
