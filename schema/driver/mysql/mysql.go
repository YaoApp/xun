package mysql

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/schema"
)

// New create new mysql blueprint instance
func New() schema.Schema {
	return &Blueprint{
		Blueprint: schema.NewBlueprint(),
	}
}

// Create Indicate that the table needs to be created.
func (blueprint *Blueprint) Create() {
	fmt.Printf("MySQL driver\n")
}
