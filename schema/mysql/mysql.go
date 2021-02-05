package mysql

import (
	"github.com/yaoapp/xun/dbal/schema"
)

// New create new mysql blueprint instance
func New(conn *schema.Connection) schema.Schema {
	return &Builder{
		Builder: schema.NewBuilder(conn),
	}
}
