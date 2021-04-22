package models

import (
	"github.com/yaoapp/xun/capsule"
	"github.com/yaoapp/xun/dbal/model"
)

// User the struct model for unit test
type User struct {
	ID       int `json:"id" x-type:"ID" x-title:"UserID" x-comment:"User id" x-validation-pattern:"^[0-9]{1,16}$"`
	Nickname string
	Address  string
	Vote     int
	Score    float64
	Status   string `x-type:"enum" x-option:"PENDING,DONE,WAITING"`
	model.Model
}

func init() {
	model.Register(&User{},
		SchemaFileContents["models/user.json"],
		SchemaFileContents["models/user.flow.json"],
	)
}

// BuildeUser create a new user instance
func BuildeUser(builder ...model.MakerFunc) User {
	user := User{}
	if len(builder) == 0 {
		builder[0] = capsule.Build
	}
	model.MakeUsing(builder[0], &user)
	return user
}

// SetAddress extend method SetAddress
func (user *User) SetAddress(address string) bool {
	user.Address = address
	return false
}
