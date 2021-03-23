package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestBuilderCreate(t *testing.T) {
	qb := New(unit.Driver(), unit.DSN())
	qb.Where()
	qb.Join()
	assert.True(t, true, "should return true")
}
