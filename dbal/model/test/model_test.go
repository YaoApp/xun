package test

import (
	"fmt"
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
	}, &user)

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
	user.Fill(row, &user)
	assert.Equal(t, "Ava", user.Get("nickname"), "The nickname should be Ava")
	assert.Equal(t, "Ava", user.Nickname, "The nickname should be Ava")
	assert.Equal(t, 99.26, user.Get("score"), "The score should be 99.26")
	assert.Equal(t, 99.26, user.Score, "The score should be 99.26")
	assert.Equal(t, nil, user.Get("not_found"), `The not_found should be nil")`)

}

func TestModelFind(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	row, err := user.Find(1)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.Equal(t, int64(1), row.Get("id"), `The return id should be 1")`)
	assert.Equal(t, 0, user.ID, `The return id should be nil")`)
	assert.Equal(t, "admin", user.Get("nickname"), `The return nickname should be admin")`)
	assert.Equal(t, "", user.Nickname, `The return nickname should be nil")`)
}

func TestModelFindEmpty(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	row, err := user.Find(2)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.True(t, row.IsEmpty(), `The return row should be empty")`)
}

func TestModelFindBind(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	_, err := user.Find(1, &user)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.Equal(t, int64(1), user.Get("id"), `The return id should be 1")`)
	assert.Equal(t, 1, user.ID, `The return id should be 1")`)
	assert.Equal(t, "admin", user.Get("nickname"), `The return nickname should be admin")`)
	assert.Equal(t, "admin", user.Nickname, `The return nickname should be admin")`)
}

func TestModelFindBindStruct(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	row := struct {
		ID       int
		Nickname string
		Bio      string
	}{}
	_, err := user.Find(1, &row)
	assert.Equal(t, nil, err, `The return error should be nil`)
	assert.Equal(t, int64(1), user.Get("id"), `The return id should be 1`)
	assert.Equal(t, 1, row.ID, `The return id should be 1`)
	assert.Equal(t, "admin", user.Get("nickname"), `The return nickname should be admin`)
	assert.Equal(t, "admin", row.Nickname, `The return nickname should be admin`)
	assert.Equal(t, "the default adminstor", row.Bio, `The return bio should be the default adminstor"`)
}

func TestModelFindBindEmpty(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	_, err := user.Find(2, &user)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.True(t, user.IsEmpty(), `The return model should be empty")`)
}

func TestModelFindSoftDeletes(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	row, err := car.Find(2)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.True(t, row.IsEmpty(), `The return model should be empty")`)
}

func TestModelSave(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)

	// Insert
	err := user.Fill(xun.R{
		"nickname":  "Ava",
		"bio":       "Yao Framework CEO",
		"user_id":   1,
		"gender":    0,
		"vote":      100,
		"score":     99.26,
		"address":   "Cecilia Chapman 711-2880 Nulla St. Mankato Mississippi 96522 (257) 563-7401",
		"status":    "DONE",
		"not_found": "something",
	}, &user).Save()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	row := user.GetQuery().Select("*").Where("nickname", "Ava").MustFirst()
	assert.Equal(t, int64(0), row.Get("gender"), `The return value should be nil")`)
	assert.Equal(t, "99.26", fmt.Sprintf("%.2f", row.Get("score")), `The return value should be nil")`)
}

func TestModelSaveUpdate(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	// Insert
	err := user.Fill(xun.R{
		"nickname":  "Ava",
		"bio":       "Yao Framework CEO",
		"user_id":   1,
		"gender":    0,
		"vote":      100,
		"score":     99.26,
		"address":   "Cecilia Chapman 711-2880 Nulla St. Mankato Mississippi 96522 (257) 563-7401",
		"status":    "DONE",
		"not_found": "something",
	}).Save(&user)
	assert.Equal(t, nil, err, `The return value should be nil")`)

	// update
	err = user.
		Set("score", 99.98, &user).
		Set("gender", 2, &user).
		Save()
	assert.Equal(t, nil, err, `The return value should be nil")`)

	row := user.GetQuery().Select("*").Where("nickname", "Ava").MustFirst()
	assert.Equal(t, int64(2), row.Get("gender"), `The return value should be nil")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", row.Get("score")), `The return value should be nil")`)
}

func TestModelSaveInsert(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	err := car.Fill(xun.R{
		"name":    "Tesla Model Y",
		"manu_id": 1,
	}).Save()
	assert.Nil(t, err, `The return value should be nil")`)

	err = car.Fill(xun.R{
		"name":    "Tesla Model Y",
		"manu_id": 1,
	}).Save()

	assert.Nil(t, err, `The return value should be nil")`)
	assert.Equal(t, int64(2), car.GetQuery().Where("name", "Tesla Model Y").MustCount(), `The return value should be 2")`)
}

func TestModelSavePrimary(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	_, err := user.Find(1, &user)
	assert.Equal(t, nil, err, `The return value should be nil")`)

	err = user.
		Set("address", "Cecilia Chapman 711-2880 Nulla St.").
		Set("score", 99.98).
		Set("bio", "Yao Framework CEO").
		Save()

	assert.Equal(t, nil, err, `The return value should be nil")`)

	_, err = user.Find(1, &user)
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.Equal(t, 1, user.ID, `The id should be 1")`)
	assert.Equal(t, "Cecilia Chapman 711-2880 Nulla St.", user.Address, `The address should be Cecilia Chapman 711-2880 Nulla St.")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", user.Score), `The score  should be 99.98")`)
}

func TestModelSaveBind(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	_, err := user.Find(1, &user)
	assert.Equal(t, nil, err, `The return value should be nil")`)

	err = user.
		Set("address", "Cecilia Chapman 711-2880 Nulla St.").
		Set("score", 99.98).
		Set("bio", "Yao Framework CEO").
		Save(&user)

	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.Equal(t, 1, user.ID, `The id should be 1")`)
	assert.Equal(t, "Cecilia Chapman 711-2880 Nulla St.", user.Address, `The address should be Cecilia Chapman 711-2880 Nulla St.")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", user.Score), `The score  should be 99.98")`)
}

func TestModelSaveBindStruct(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.MakeUser(modelTestMaker)
	_, err := user.Find(1, &user)
	assert.Equal(t, nil, err, `The return value should be nil")`)
	row := struct {
		ID      int
		Score   float64
		Address string
		Bio     string
	}{}

	err = user.
		Set("address", "Cecilia Chapman 711-2880 Nulla St.").
		Set("score", 99.98).
		Set("bio", "Yao Framework CEO").
		Save(&row)

	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.Equal(t, 1, row.ID, `The id should be 1")`)
	assert.Equal(t, "Yao Framework CEO", row.Bio, `The bio should be Yao Framework CEO")`)
	assert.Equal(t, "Cecilia Chapman 711-2880 Nulla St.", row.Address, `The address should be Cecilia Chapman 711-2880 Nulla St.")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", row.Score), `The score  should be 99.98")`)
}
