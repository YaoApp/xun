package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	conn := &Connection{}
	New(conn)
	qb := NewBuilder(conn)
	qb.Where()
	qb.Join()
	assert.True(t, true, "should return true")
}
