package capsule

import (
	"github.com/yaoapp/xun/dbal/schema"
)

// Get a schema builder instance.
func newSchema(driver string, conn *schema.Connection) schema.Schema {
	return schema.Use(conn)
}
