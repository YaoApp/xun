package test

import (
	"testing"

	"github.com/yaoapp/xun/dbal/model"
	"github.com/yaoapp/xun/dbal/model/test/models"
	"github.com/yaoapp/xun/dbal/query"
	"github.com/yaoapp/xun/unit"
)

var testQueryBuilder query.Query

func makeModelTest(v interface{}) *model.Model {
	defer unit.Catch()
	unit.SetLogger()
	if testQueryBuilder != nil {
		return model.Make(testQueryBuilder, v)
	}
	testQueryBuilder = query.New(unit.Driver(), unit.DSN())
	return model.Make(testQueryBuilder, v)
}

func TestMakeByJsonSchema(t *testing.T) {
	defer unit.Catch()
	member := makeModelTest(models.SchemaFileContents["models/member.json"])
	member.Search()
	member.Find()
}

func TestMakeByStruct(t *testing.T) {
	defer unit.Catch()
	user := &models.User{}
	makeModelTest(user)
	user.SetAddress("hello")
	user.Search()
	user.Find()
}
