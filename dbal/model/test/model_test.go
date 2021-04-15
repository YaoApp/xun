package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal/model"
	"github.com/yaoapp/xun/dbal/model/test/models"
	"github.com/yaoapp/xun/utils"
)

func TestModelColumns(t *testing.T) {
	registerModelsForTest()
	user := models.MakeUser(modelTestMaker)
	columns := user.Columns()
	assert.Equal(t, 8, len(columns), "The return columns count of model should be 8")

	member := model.MakeUsing(modelTestMaker, "models/member")
	columns = member.Columns()
	assert.Equal(t, 6, len(columns), "The return columns count of model should be 6")
}

func TestModelSearchable(t *testing.T) {
	registerModelsForTest()
	user := models.MakeUser(modelTestMaker)
	columns := user.Searchable()
	searchable := []string{"status", "gender", "id", "nickname", "score"}
	for _, column := range searchable {
		assert.True(t, utils.StringHave(columns, column), "The return searchable columns should have %s", column)
	}

	member := model.MakeUsing(modelTestMaker, "models/member")
	columns = member.Searchable()
	searchable = []string{"id", "user_id"}
	for _, column := range searchable {
		assert.True(t, utils.StringHave(columns, column), "The return searchable columns should have %s", column)
	}
}

func TestModelPrimaryKeys(t *testing.T) {
	registerModelsForTest()
	user := models.MakeUser(modelTestMaker)
	columns := user.PrimaryKeys()
	primaryKeys := []string{"id"}
	for _, column := range primaryKeys {
		assert.True(t, utils.StringHave(columns, column), "The return primary key  should have %s", column)
	}

	member := model.MakeUsing(modelTestMaker, "models/member")
	columns = member.PrimaryKeys()
	primaryKeys = []string{"id"}
	for _, column := range primaryKeys {
		assert.True(t, utils.StringHave(columns, column), "The return  primary key should have %s", column)
	}
}

func TestModelPrimary(t *testing.T) {
	registerModelsForTest()
	user := models.MakeUser(modelTestMaker)
	assert.Equal(t, "id", user.Primary(), "The return primary key should be id")

	member := model.MakeUsing(modelTestMaker, "models/member")
	assert.Equal(t, "id", member.Primary(), "The return primary key should be id")
}
