package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	goSQL "github.com/yaoapp/xun/grammar/sql"
)

func newTestPostgres() Postgres {
	return Postgres{
		SQL: goSQL.NewSQL(&Quoter{}),
	}
}

func TestWhereJsoncontainsPG(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `"tags"::jsonb @> $1`, result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsPGNot(t *testing.T) {
	pg := newTestPostgres()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `not "tags"::jsonb @> $1`, result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsPGOffset(t *testing.T) {
	pg := newTestPostgres()
	offset := 2
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"test"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := pg.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, `"tags"::jsonb @> $3`, result)
	assert.Equal(t, 3, offset)
}

func TestWhereJsoncontainsPGMixed(t *testing.T) {
	pg := newTestPostgres()
	offset := 0

	where1 := dbal.Where{
		Type:     "basic",
		Column:   "name",
		Operator: "=",
		Value:    "test",
		Boolean:  "and",
		Not:      false,
		Offset:   1,
	}
	result1 := pg.WhereBasic(&dbal.Query{}, where1, &offset)
	assert.Equal(t, `"name" = $1`, result1)
	assert.Equal(t, 1, offset)

	where2 := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   `"admin"`,
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}
	result2 := pg.WhereJsoncontains(&dbal.Query{}, where2, &offset)
	assert.Equal(t, `"tags"::jsonb @> $2`, result2)
	assert.Equal(t, 2, offset)
}
