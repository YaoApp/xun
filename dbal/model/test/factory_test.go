package test

import (
	"fmt"
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

func getSchema() schema.Schema {
	unit.SetLogger()
	if testSchemaBuilder != nil {
		return testSchemaBuilder
	}
	testSchemaBuilder = schema.New(unit.Driver(), unit.DSN())
	return testSchemaBuilder
}

func getQuery() query.Query {
	unit.SetLogger()
	if testQueryBuilder != nil {
		return testQueryBuilder
	}
	testQueryBuilder = query.New(unit.Driver(), unit.DSN())
	return testQueryBuilder
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
	assert.Equal(t, 10, len(user.GetAttributeNames()), "The return attribute names count of model should be 10")
	assert.Equal(t, 10, len(user.GetAttributes()), "The return attribute names count of model should be 10")

	status := user.GetAttr("status")
	assert.Equal(t, 3, len(status.Column.Option), "The status column option count should be 3")
	assert.Equal(t, "WAITING", status.Column.Default, "The status column default value should be WAITING")
}

func TestFactoryRegisterError(t *testing.T) {
	assert.PanicsWithError(t, "The type kind (string) can't be register, have 0 arguments", func() {
		model.MakeUsing(modelTestMaker, "models/car")
	})
}

func TestFactoryRegister(t *testing.T) {
	car := model.MakeUsing(modelTestMaker, "models/car", models.SchemaFileContents["models/car.json"])
	assert.Equal(t, "models.car", car.GetFullname(), "The return fullname of model should be models.car")
	assert.Equal(t, 5, len(car.GetAttributeNames()), "The return attribute names count of model should be 5")
	assert.Equal(t, 5, len(car.GetAttributes()), "The return attribute names count of model should be 5")

	manu := model.MakeUsing(modelTestMaker, "models/manu", models.SchemaFileContents["models/manu.json"])
	assert.Equal(t, "models.manu", manu.GetFullname(), "The return fullname of model should be models.car")
	assert.Equal(t, 4, len(manu.GetAttributeNames()), "The return attribute names count of model should be 4")
	assert.Equal(t, 4, len(manu.GetAttributes()), "The return attribute names count of model should be 4")

	userCar := model.MakeUsing(modelTestMaker, "models/user_car", models.SchemaFileContents["models/user_car.json"])
	assert.Equal(t, "models.user_car", userCar.GetFullname(), "The return fullname of model should be models.user_car")
	assert.Equal(t, 3, len(userCar.GetAttributeNames()), "The return attribute names count of model should be 3")
	assert.Equal(t, 3, len(userCar.GetAttributes()), "The return attribute names count of model should be 3")
}

func TestFactoryClass(t *testing.T) {
	registerModelsForTest()
	models := []string{"user", "member", "manu", "car", "user_car", "null"}
	for _, name := range models {
		name = fmt.Sprintf("models.%s", name)
		factory := model.Class(name)
		model := model.GetModel(factory.Model)
		assert.True(t, name == model.GetFullname(), "the  model.GetFullname() shoud be %s ", name)
	}
}

func TestFactoryClassError(t *testing.T) {
	registerModelsForTest()
	assert.PanicsWithError(t, "The model (notfound) doesn't register", func() {
		model.Class("notfound")
	})
}

func TestFactoryGetModelError(t *testing.T) {
	registerModelsForTest()
	assert.PanicsWithError(t, "v is (*string) not a model", func() {
		v := "notfound"
		model.GetModel(&v)
	})
}

func TestFactorySetModelError(t *testing.T) {
	new := model.Model{}
	assert.PanicsWithError(t, "v is (*string) not a model", func() {
		v := "notfound"
		model.SetModel(&v, new)
	})
}

func TestFactoryNewError(t *testing.T) {
	assert.PanicsWithError(t, "The model type (string) must be a pointer", func() {
		v := "notfound"
		model.Class("models.user").New(v)
	})
}

func TestFactoryMigrate(t *testing.T) {
	registerModelsForTest()
	sch := getSchema()
	qb := getQuery()
	models := []string{"user", "member", "manu", "car", "user_car", "null"}
	for _, name := range models {
		name = fmt.Sprintf("models.%s", name)
		err := model.Class(name).Migrate(sch, qb, true)
		assert.Nil(t, err, "create %s the return value should be nil", name)
	}
}

func TestFactoryMigrateError(t *testing.T) {
	registerModelsForTest()
	assert.PanicsWithError(t, `This feature does not support it yet. It working when the first parameter refresh is true.(model.Class("user").Migrate(schema, true))`, func() {
		sch := getSchema()
		qb := getQuery()
		model.Class("models.user").Migrate(sch, qb)
	})
}

func TestFactoryMethods(t *testing.T) {
	registerModelsForTest()
	methods := model.Class("models.user").Methods()
	assert.Greater(t, len(methods), 1, "The return attribute names count of model should be greater 1")

	methods = model.Class("models.member").Methods()
	assert.Greater(t, len(methods), 1, "The return attribute names count of model should be greater 1")
}

func TestFactoryMethodsCached(t *testing.T) {
	registerModelsForTest()
	methods := model.Class("models.user").Methods()
	assert.Greater(t, len(methods), 1, "The return attribute names count of model should be greater 1")

	methods = model.Class("models.member").Methods()
	assert.Greater(t, len(methods), 1, "The return attribute names count of model should be greater 1")
}

// Utils ...

func registerModelsForTest() {
	modelNames := []string{"member", "manu", "car", "user_car", "null"}
	for _, name := range modelNames {
		args := []interface{}{}
		if schema, has := models.SchemaFileContents[fmt.Sprintf("models/%s.json", name)]; has {
			args = append(args, schema)
		}
		if flow, has := models.SchemaFileContents[fmt.Sprintf("models/%s.flow.json", name)]; has {
			args = append(args, flow)
		}
		if len(args) < 1 {
			continue
		}
		model.MakeUsing(modelTestMaker, fmt.Sprintf("models/%s", name), args...)
	}
}
