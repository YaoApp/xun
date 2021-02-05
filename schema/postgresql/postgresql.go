package postgresql

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/schema"
)

// New create new mysql blueprint instance
func New(conn *schema.Connection) schema.Schema {
	return &Builder{
		Builder: schema.NewBuilder(conn),
	}
}

// Create Indicate that the table needs to be created.
func (buider *Builder) Create() {
	fmt.Printf("\nCreate postreSQL: \n===\n%#v\n===\n", buider.Conn)
}
