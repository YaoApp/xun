package capsule

import (
	"testing"
)

func TestNew(t *testing.T) {
	db := New()
	db.Schema()
}
