package mysql

import (
	"testing"

	"github.com/yaoapp/xun/unit"
)

func TestCreate(t *testing.T) {
	db := unit.Use("mysql")
	schema := db.Schema()
	schema.Create()
	schema.Drop()
}
