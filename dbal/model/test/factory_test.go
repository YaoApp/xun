package test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal/model"
	"github.com/yaoapp/xun/dbal/model/test/models"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

var testQueryBuilder query.Query
var testSchemaBuilder schema.Schema

func modelTestMaker(v interface{}, flow ...interface{}) *model.Model {
	unit.SetLogger()
	if testQueryBuilder != nil && testSchemaBuilder != nil {
		return model.Make(testQueryBuilder, testSchemaBuilder, v, flow...)
	}
	testQueryBuilder = query.New(unit.Driver(), unit.DSN())
	testSchemaBuilder = schema.New(unit.Driver(), unit.DSN())
	return model.Make(testQueryBuilder, testSchemaBuilder, v, flow...)
}

func TestMakeBySchema(t *testing.T) {
	member := model.MakeUsing(modelTestMaker, "models/member",
		models.SchemaFileContents["models/member.json"],
		models.SchemaFileContents["models/member.flow.json"],
	)
	assert.Equal(t, "*model.Model", reflect.TypeOf(member).String(), "The return type of member should be *model.Model")
	assert.Equal(t, "models", member.GetNamespace(), "The return namespace of model should be models.member")
	assert.Equal(t, "member", member.GetName(), "The return name of model should be member")
	assert.Equal(t, "models.member", member.GetFullname(), "The return fullname of model should be models.member")

}

func TestMakeBySchemaRegistered(t *testing.T) {
	member := model.MakeUsing(modelTestMaker, "models/member")
	assert.Equal(t, "*model.Model", reflect.TypeOf(member).String(), "The return type of member should be *model.Model")
	assert.Equal(t, "models", member.GetNamespace(), "The return namespace of model should be models.member")
	assert.Equal(t, "member", member.GetName(), "The return name of model should be member")
	assert.Equal(t, "models.member", member.GetFullname(), "The return fullname of model should be models.member")
}

func TestMakeByStruct(t *testing.T) {
	user := models.User{}
	modelTestMaker(&user)
	model.MakeUsing(modelTestMaker, &user)
	user.SetAddress("hello new address")
	assert.Equal(t, "models.User", reflect.TypeOf(user).String(), "The return type of member should be models.User")
	assert.Equal(t, "hello new address", user.Address, "The user address should be hello new address")
	assert.Equal(t, "models", user.GetNamespace(), "The return namespace of model should be models.member")
	assert.Equal(t, "user", user.GetName(), "The return name of model should be member")
	assert.Equal(t, "models.user", user.GetFullname(), "The return fullname of model should be models.member")
}

func TestMakeByStructStyle2(t *testing.T) {
	user := models.MakeUser(modelTestMaker)
	user.SetAddress("hello new address style2")
	assert.Equal(t, "models.User", reflect.TypeOf(user).String(), "The return type of member should be models.User")
	assert.Equal(t, "hello new address style2", user.Address, "The user address should be hello new address")
	assert.Equal(t, "models", user.GetNamespace(), "The return namespace of model should be models.member")
	assert.Equal(t, "user", user.GetName(), "The return name of model should be member")
	assert.Equal(t, "models.user", user.GetFullname(), "The return fullname of model should be models.member")
}

func TestFactoryRegisterError(t *testing.T) {
	assert.PanicsWithError(t, "The type kind (string) can't be register, have 0 arguments", func() {
		model.MakeUsing(modelTestMaker, "models/car")
	})
}

func TestFactoryRegister(t *testing.T) {
	car := model.MakeUsing(modelTestMaker, "models/car", models.SchemaFileContents["models/car.json"])
	assert.Equal(t, "models.car", car.GetFullname(), "The return fullname of model should be models.car")

	manu := model.MakeUsing(modelTestMaker, "models/manu", models.SchemaFileContents["models/manu.json"])
	assert.Equal(t, "models.manu", manu.GetFullname(), "The return fullname of model should be models.car")

	userCar := model.MakeUsing(modelTestMaker, "models/user_car", models.SchemaFileContents["models/user_car.json"])
	assert.Equal(t, "models.user_car", userCar.GetFullname(), "The return fullname of model should be models.user_car")
}
