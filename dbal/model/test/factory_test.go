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
	defer unit.Catch()
	unit.SetLogger()
	if testQueryBuilder != nil && testSchemaBuilder != nil {
		return model.Make(testQueryBuilder, testSchemaBuilder, v, flow...)
	}
	testQueryBuilder = query.New(unit.Driver(), unit.DSN())
	testSchemaBuilder = schema.New(unit.Driver(), unit.DSN())
	return model.Make(testQueryBuilder, testSchemaBuilder, v, flow...)
}

func TestMakeBySchema(t *testing.T) {
	defer unit.Catch()
	member := model.MakeUsing(modelTestMaker,
		models.SchemaFileContents["models/member.json"],
		models.SchemaFileContents["models/member.flow.json"],
	)
	assert.Equal(t, "*model.Model", reflect.TypeOf(member).String(), "The return type of member should be *model.Model")
}

func TestMakeBySchemaRegistered(t *testing.T) {
	defer unit.Catch()
	member := model.MakeUsing(modelTestMaker,
		models.SchemaFileContents["models/member.json"],
		models.SchemaFileContents["models/member.flow.json"],
	)
	assert.Equal(t, "*model.Model", reflect.TypeOf(member).String(), "The return type of member should be *model.Model")
}

func TestMakeByStruct(t *testing.T) {
	defer unit.Catch()
	user := models.User{}
	modelTestMaker(&user)
	model.MakeUsing(modelTestMaker, &user)
	user.SetAddress("hello new address")
	assert.Equal(t, "models.User", reflect.TypeOf(user).String(), "The return type of member should be models.User")
	assert.Equal(t, "hello new address", user.Address, "The user address should be hello new address")
}

func TestMakeByStructStyle2(t *testing.T) {
	defer unit.Catch()
	user := models.MakeUser(modelTestMaker)
	user.SetAddress("hello new address style2")
	assert.Equal(t, "models.User", reflect.TypeOf(user).String(), "The return type of member should be models.User")
	assert.Equal(t, "hello new address style2", user.Address, "The user address should be hello new address")
}
