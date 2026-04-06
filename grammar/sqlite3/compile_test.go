package sqlite3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	goSQL "github.com/yaoapp/xun/grammar/sql"
)

func newTestSQLite3() SQLite3 {
	return SQLite3{
		SQL: goSQL.NewSQL(&goSQL.Quoter{}),
	}
}

func TestWhereJsoncontainsSQLite(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%admin%",
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "`tags` like ?", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsSQLiteNot(t *testing.T) {
	g := newTestSQLite3()
	offset := 0
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%admin%",
		Boolean: "and",
		Not:     true,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "not `tags` like ?", result)
	assert.Equal(t, 1, offset)
}

func TestWhereJsoncontainsSQLiteOffset(t *testing.T) {
	g := newTestSQLite3()
	offset := 5
	where := dbal.Where{
		Type:    "jsoncontains",
		Column:  "tags",
		Value:   "%test%",
		Boolean: "and",
		Not:     false,
		Offset:  1,
	}

	result := g.WhereJsoncontains(&dbal.Query{}, where, &offset)
	assert.Equal(t, "`tags` like ?", result)
	assert.Equal(t, 6, offset)
}
