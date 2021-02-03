package oracle

import (
	"fmt"

	"github.com/yaoapp/xun/dbal"
)

// New create new mysql blueprint instance
func New() dbal.Schema {
	return &Blueprint{
		Blueprint: dbal.Blueprint{},
	}
}

// Create Indicate that the table needs to be created.
func (blueprint *Blueprint) Create() {
	fmt.Printf("Oracle driver\n")
}
