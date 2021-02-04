package capsule

import (
	"testing"
)

func TestNew(t *testing.T) {
	db := New()
	schema := db.Schema()
	schema.Drop()

	query := db.Table()
	query.Where()
}
