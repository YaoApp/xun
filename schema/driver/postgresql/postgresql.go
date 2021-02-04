package postgresql

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/schema"
)

// New create new mysql blueprint instance
func New(conn *schema.Connection) schema.Schema {
	return &Blueprint{
		Blueprint: schema.NewBlueprint(conn),
	}
}

// Create Indicate that the table needs to be created.
func (blueprint *Blueprint) Create() {
	fmt.Printf("postgreSQL driver\n")
}
