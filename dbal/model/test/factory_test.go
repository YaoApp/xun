package test

import (
	"testing"

	"github.com/yaoapp/xun/dbal/model"
	"github.com/yaoapp/xun/dbal/model/test/models"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

var testQueryBuilder query.Query
var testSchemaBuilder schema.Schema

func makeModelTest(v interface{}, flow ...interface{}) *model.Model {
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
	member := makeModelTest(
		models.SchemaFileContents["models/member.json"],
		models.SchemaFileContents["models/member.flow.json"],
	)
	member.Search()
	member.Find()
}

func TestMakeBySchemaRegistered(t *testing.T) {
	defer unit.Catch()
	member := makeModelTest(
		models.SchemaFileContents["models/member.json"],
		models.SchemaFileContents["models/member.flow.json"],
	)
	member.Search()
	member.Find()
}

func TestMakeByStruct(t *testing.T) {
	defer unit.Catch()
	user := models.User{}
	makeModelTest(&user)
	user.SetAddress("hello new address")
	u := user.Find()
	utils.Println(u)
	user.Search()
}
