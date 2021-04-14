package models

import (
	"github.com/yaoapp/xun/capsule"
	"github.com/yaoapp/xun/dbal/model"
)

// User the struct model for unit test
type User struct {
	ID       int `json:"id" x-title:"UserID" x-comment:"The user id" x-validation-pattern:"^[0-9]{1,16}$"`
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

// MakeUser create a new user instance
func MakeUser(maker ...model.MakerFunc) User {
	user := User{}
	if len(maker) == 0 {
		maker[0] = capsule.Make
	}
	model.MakeUsing(maker[0], &user)
	return user
}

// SetAddress extend method SetAddress
func (user *User) SetAddress(address string) bool {
	user.Address = address
	return false
}

// Find user fild
func (user *User) Find() *User {
	return user.Model.Find(user).(*User)
}
