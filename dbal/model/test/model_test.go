package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
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

func TestModelFillSchemaXunR(t *testing.T) {
	registerModelsForTest()
	member := model.MakeUsing(modelTestMaker, "models/member")
	member.Fill(xun.R{
		"user_id":    1,
		"name":       "Emma",
		"score":      99.26,
		"level":      "gold",
		"expired_at": dbal.Raw("NOW()"),
		"not_found":  "something",
	})
	assert.Equal(t, 1, member.Get("user_id"), "The user_id should be 1")
	assert.Equal(t, "Emma", member.Get("name"), "The name should be Emma")
	assert.Equal(t, 99.26, member.Get("score"), "The score should be 99.26")
	assert.Equal(t, "gold", member.Get("level"), "The level should be gold")
	assert.Equal(t, dbal.Raw("NOW()"), member.Get("expired_at"), `The expired_at should be dbal.Raw("NOW()")`)
	assert.Equal(t, nil, member.Get("not_found"), `The not_found should be nil")`)
}

func TestModelFillStructXunR(t *testing.T) {
	registerModelsForTest()
	user := models.MakeUser(modelTestMaker)
	user.Fill(xun.R{
		"nickname":  "Ava",
		"bio":       "Yao Framework CEO",
		"gender":    0,
		"vote":      100,
		"score":     99.26,
		"address":   "Cecilia Chapman 711-2880 Nulla St. Mankato Mississippi 96522 (257) 563-7401",
		"status":    "DONE",
		"not_found": "something",
	})
	assert.Equal(t, "Ava", user.Get("nickname"), "The nickname should be Ava")
	assert.Equal(t, "Ava", user.Nickname, "The nickname should be Ava")
	assert.Equal(t, "Yao Framework CEO", user.Get("bio"), "The nickname should be Yao Framework CEO")
	assert.Equal(t, 99.26, user.Get("score"), "The score should be 99.26")
	assert.Equal(t, 99.26, user.Score, "The score should be 99.26")
	assert.Equal(t, nil, user.Get("not_found"), `The not_found should be nil")`)
}

func TestModelFillStructUser(t *testing.T) {
	registerModelsForTest()
	row := models.User{
		Nickname: "Ava",
		Vote:     100,
		Score:    99.26,
		Status:   "DONE",
	}
	user := models.MakeUser(modelTestMaker)
	user.Fill(row)
	assert.Equal(t, "Ava", user.Get("nickname"), "The nickname should be Ava")
	assert.Equal(t, "Ava", user.Nickname, "The nickname should be Ava")
	assert.Equal(t, 99.26, user.Get("score"), "The score should be 99.26")
	assert.Equal(t, 99.26, user.Score, "The score should be 99.26")
	assert.Equal(t, nil, user.Get("not_found"), `The not_found should be nil")`)
}
