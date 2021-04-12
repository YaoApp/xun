package models

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/model"
)

// User the struct model for unit test
type User struct {
	ID      int `json:"id" x-title:"UserID" x-comment:"The user id" x-validation-pattern:"^[0-9]{1,16}$"`
	Name    string
	Address string
	Vote    int
	Score   float64
	Status  string `x-type:"enum" x-option:"PENDING,DONE,WAITING"`
	model.Model
}

func init() {
	model.Register(&User{}, SchemaFileContents["models/user.json"])
}

// SetAddress extend method SetAddress
func (user *User) SetAddress(address string) bool {
	return false
}

// Find user fild
func (user *User) Find() {
	fmt.Println("User find")
}
