package schema

import (
	"fmt"
)

// New create new schema buider interface
func New() Schema {
	blueprint := NewBlueprint()
	return &blueprint
}

// NewBlueprint create new schema buider blueprint
func NewBlueprint() Blueprint {
	return Blueprint{}
}

// Create Indicate that the table needs to be created.
func (blueprint *Blueprint) Create() {}

// Drop Indicate that the table should be dropped.
func (blueprint *Blueprint) Drop() {
	fmt.Printf("Drop DBAL\n")
}

// DropIfExists Indicate that the table should be dropped if it exists.
func (blueprint *Blueprint) DropIfExists() {}

// Rename the table to a given name.
func (blueprint *Blueprint) Rename() {}

// Primary Specify the primary key(s) for the table.
func (blueprint *Blueprint) Primary() {}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (blueprint *Blueprint) BigInteger() {}

// String Create a new string column on the table.
func (blueprint *Blueprint) String() {}
