package model

import "fmt"

// Register register the grammar driver
func Register(v interface{}, schema ...interface{}) error {
	return nil
}

// Search search by given params
func (model *Model) Search() {
}

// Find find by primary key
func (model *Model) Find() {
	fmt.Println("Model find")
}
