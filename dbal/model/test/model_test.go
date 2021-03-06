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
	user := models.BuildeUser(modelTestMaker)
	columns := user.Columns()
	assert.Equal(t, 8, len(columns), "The return columns count of model should be 8")

	member := model.MakeUsing(modelTestMaker, "models/member")
	columns = member.Columns()
	assert.Equal(t, 6, len(columns), "The return columns count of model should be 6")
}

func TestModelSearchable(t *testing.T) {
	registerModelsForTest()
	user := models.BuildeUser(modelTestMaker)
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
	user := models.BuildeUser(modelTestMaker)
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
	user := models.BuildeUser(modelTestMaker)
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
	assert.Equal(t, 1, member.Value("user_id"), "The user_id should be 1")
	assert.Equal(t, "Emma", member.Value("name"), "The name should be Emma")
	assert.Equal(t, 99.26, member.Value("score"), "The score should be 99.26")
	assert.Equal(t, "gold", member.Value("level"), "The level should be gold")
	assert.Equal(t, dbal.Raw("NOW()"), member.Value("expired_at"), `The expired_at should be dbal.Raw("NOW()")`)
	assert.Equal(t, nil, member.Value("not_found"), `The not_found should be nil")`)
}

func TestModelFillStructXunR(t *testing.T) {
	registerModelsForTest()
	user := models.BuildeUser(modelTestMaker)
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

	assert.Equal(t, "Ava", user.Value("nickname"), "The nickname should be Ava")
	assert.Equal(t, "Ava", user.Nickname, "The nickname should be Ava")
	assert.Equal(t, "Yao Framework CEO", user.Value("bio"), "The nickname should be Yao Framework CEO")
	assert.Equal(t, 99.26, user.Value("score"), "The score should be 99.26")
	assert.Equal(t, 99.26, user.Score, "The score should be 99.26")
	assert.Equal(t, nil, user.Value("not_found"), `The not_found should be nil")`)
}

func TestModelFillStructUser(t *testing.T) {
	registerModelsForTest()
	row := models.User{
		Nickname: "Ava",
		Vote:     100,
		Score:    99.26,
		Status:   "DONE",
	}
	user := models.BuildeUser(modelTestMaker)
	user.Fill(row, &user)
	assert.Equal(t, "Ava", user.Value("nickname"), "The nickname should be Ava")
	assert.Equal(t, "Ava", user.Nickname, "The nickname should be Ava")
	assert.Equal(t, 99.26, user.Value("score"), "The score should be 99.26")
	assert.Equal(t, 99.26, user.Score, "The score should be 99.26")
	assert.Equal(t, nil, user.Value("not_found"), `The not_found should be nil")`)

}

func TestModelFind(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	row, err := user.Find(1)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.Equal(t, int64(1), row.Value("id"), `The return id should be 1")`)
	assert.Equal(t, 0, user.ID, `The return id should be nil")`)
	assert.Equal(t, "admin", user.Value("nickname"), `The return nickname should be admin")`)
	assert.Equal(t, "", user.Nickname, `The return nickname should be nil")`)
}

func TestModelFindEmpty(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	row, err := user.Find(2)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.True(t, row.IsEmpty(), `The return row should be empty")`)
}

func TestModelFindBind(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	_, err := user.Find(1, &user)
	assert.Equal(t, nil, err, `The return error should be nil")`)
	assert.Equal(t, int64(1), user.Value("id"), `The return id should be 1")`)
	assert.Equal(t, 1, user.ID, `The return id should be 1")`)
	assert.Equal(t, "admin", user.Value("nickname"), `The return nickname should be admin")`)
	assert.Equal(t, "admin", user.Nickname, `The return nickname should be admin")`)
}

func TestModelFindBindStruct(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	row := struct {
		ID       int
		Nickname string
		Bio      string
	}{}
	_, err := user.Find(1, &row)
	assert.Equal(t, nil, err, `The return error should be nil`)
	assert.Equal(t, int64(1), user.Value("id"), `The return id should be 1`)
	assert.Equal(t, 1, row.ID, `The return id should be 1`)
	assert.Equal(t, "admin", user.Value("nickname"), `The return nickname should be admin`)
	assert.Equal(t, "admin", row.Nickname, `The return nickname should be admin`)
	assert.Equal(t, "the default adminstor", row.Bio, `The return bio should be the default adminstor"`)
}

func TestModelFindBindEmpty(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
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
	user := models.BuildeUser(modelTestMaker)

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
	row := user.Select("*").Where("nickname", "Ava").MustFirst()
	assert.Equal(t, int64(0), row.Value("gender"), `The return value should be nil")`)
	assert.Equal(t, "99.26", fmt.Sprintf("%.2f", row.Value("score")), `The return value should be nil")`)
}

func TestModelSaveUpdate(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
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

	row := user.Select("*").Where("nickname", "Ava").MustFirst()
	assert.Equal(t, int64(2), row.Value("gender"), `The return value should be nil")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", row.Value("score")), `The return value should be nil")`)
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
	assert.Equal(t, int64(2), car.Where("name", "Tesla Model Y").MustCount(), `The return value should be 2")`)
}

func TestModelSaveUniqueKeys(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	err := user.Fill(xun.R{
		"nickname": "admin",
		"vote":     100,
		"score":    29.99,
		"bio":      "Yao Framework CEO",
		"address":  "Cecilia Chapman 711-2880 Nulla St.",
	}).Save()

	assert.Nil(t, err, `The return value should be nil")`)
	_, err = user.Find(1, &user)
	assert.Nil(t, err, `The return value should be nil")`)
	assert.Equal(t, 1, user.ID, `The return value should1`)
	assert.Equal(t, "Cecilia Chapman 711-2880 Nulla St.", user.Address, `The return value should be Cecilia Chapman 711-2880 Nulla St.`)
}

func TestModelSavePrimary(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
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
	user := models.BuildeUser(modelTestMaker)
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
	user := models.BuildeUser(modelTestMaker)
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

	assert.Equal(t, nil, err, `The return value should be nil`)
	assert.Equal(t, 1, row.ID, `The id should be 1")`)
	assert.Equal(t, "Yao Framework CEO", row.Bio, `The bio should be Yao Framework CEO")`)
	assert.Equal(t, "Cecilia Chapman 711-2880 Nulla St.", row.Address, `The address should be Cecilia Chapman 711-2880 Nulla St.")`)
	assert.Equal(t, "99.98", fmt.Sprintf("%.2f", row.Score), `The score  should be 99.98")`)
}

func TestModelDestory(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	assert.False(t, user.MustFind(1).IsEmpty(), `The user table should have 1 row`)

	err := user.MustFind(1).Destroy()
	assert.Equal(t, nil, err, `The return value should be nil"`)
	assert.True(t, user.MustFind(1).IsEmpty(), `The return value should be true"`)
}

func TestModelDestoryByID(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	assert.False(t, user.MustFind(1).IsEmpty(), `The user table should have 1 row`)

	err := user.Destroy(1)
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, user.MustFind(1).IsEmpty(), `The return value should be true"`)
}

func TestModelDestoryByIDs(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	err := user.Fill(xun.R{
		"nickname": "Ava",
		"vote":     100,
		"score":    29.99,
		"bio":      "Yao Framework CEO",
		"address":  "Cecilia Chapman 711-2880 Nulla St.",
	}).Save()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.False(t, user.MustFind(1).IsEmpty(), `The user table id = 1 should not empty `)
	assert.False(t, user.MustFind(2).IsEmpty(), `The user table id = 2 should not empty`)

	err = user.Destroy([]int{1, 2})
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, user.MustFind(1).IsEmpty(), `The return value should be true"`)
	assert.True(t, user.MustFind(2).IsEmpty(), `The return value should be true"`)
}

func TestModelDestoryByIDsStyle2(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	err := user.Fill(xun.R{
		"nickname": "Ava",
		"vote":     100,
		"score":    29.99,
		"bio":      "Yao Framework CEO",
		"address":  "Cecilia Chapman 711-2880 Nulla St.",
	}).Save()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.False(t, user.MustFind(1).IsEmpty(), `The user table id = 1 should not empty `)
	assert.False(t, user.MustFind(2).IsEmpty(), `The user table id = 2 should not empty`)

	err = user.Destroy(1, 2)
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, user.MustFind(1).IsEmpty(), `The return value should be true"`)
	assert.True(t, user.MustFind(2).IsEmpty(), `The return value should be true"`)
}

func TestModelDestorySoftDeletes(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	assert.False(t, car.MustFind(1).IsEmpty(), `The car table should have 1 row`)

	err := car.MustFind(1).Destroy()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, car.MustFind(1).IsEmpty(), `The return value should be true"`)

	row := car.Reset().WithTrashed().Select("*").Where("id", 1).MustFirst()
	assert.Equal(t, int64(1), row.Value("id"), `The return value should be 1")`)
	assert.NotNil(t, row.Value("deleted_at"), `The return value should be datetime")`)
}

func TestModelWithTrashed(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	assert.False(t, car.MustFind(1).IsEmpty(), `The car table should have 1 row`)

	err := car.MustFind(1).Destroy()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, car.MustFind(1).IsEmpty(), `The return value should be true"`)
	assert.False(t, car.WithTrashed().MustFind(1).IsEmpty(), `The return value should be false"`)
}

func TestModelOnlyTrashed(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	assert.False(t, car.MustFind(1).IsEmpty(), `The car table should have 1 row`)
	assert.True(t, car.OnlyTrashed().MustFind(1).IsEmpty(), `The return value should be false"`)

	err := car.MustFind(1).Destroy()
	assert.Equal(t, nil, err, `The return value should be nil")`)
	assert.True(t, car.MustFind(1).IsEmpty(), `The return value should be true"`)
	assert.False(t, car.OnlyTrashed().MustFind(1).IsEmpty(), `The return value should be false"`)
}

func TestModelQuery(t *testing.T) {
	TestFactoryMigrate(t)
	user := models.BuildeUser(modelTestMaker)
	rows := user.MustGet()
	assert.Equal(t, 1, len(rows), `The return value should be 1")`)
	if len(rows) == 1 {
		assert.Equal(t, int64(1), rows[0].Value("id"), `The return value should be 1")`)
	}
}

func TestModelQuerySoftdeletes(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	rows := car.MustGet()
	assert.Equal(t, 6, len(rows), `The return value should be 6")`)
	if len(rows) == 6 {
		for i := range rows {
			assert.Nil(t, rows[i].Value("deleted_at"), `The return value should be datetime")`)
		}
	}
}

func TestModelQuerySoftdeletesWithTrashed(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	rows := car.MustGet()
	assert.Equal(t, 6, len(rows), `The return value should be 6")`)
	if len(rows) == 6 {
		for i := range rows {
			assert.Nil(t, rows[i].Value("deleted_at"), `The return value should be datetime")`)
		}
	}

	car.Reset()
	rows = car.WithTrashed().MustGet()
	assert.Equal(t, 7, len(rows), `The return value should be 2")`)
	if len(rows) == 7 {
		assert.Nil(t, rows[0].Value("deleted_at"), `The return value should be nil")`)
		assert.NotNil(t, rows[1].Value("deleted_at"), `The return value should be datetime")`)
	}
}

func TestModelQuerySoftdeletesOnlyTrashed(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	rows := car.OnlyTrashed().MustGet()
	assert.Equal(t, 1, len(rows), `The return value should be 1")`)
	if len(rows) == 1 {
		assert.NotNil(t, rows[0].Value("deleted_at"), `The return value should be datetime")`)
	}
}

func TestModelWithHasOne(t *testing.T) {
	TestFactoryMigrate(t)
	car := model.MakeUsing(modelTestMaker, "models/car")
	rows := car.With("manu").MustGet()
	assert.Equal(t, 6, len(rows), `The return value should be 1")`)
	if len(rows) == 6 {
		for i := range rows {
			assert.Equal(t, rows[i].Get("manu_id"), rows[i].Get("manu.id"), `The return value should be 1")`)
			assert.True(t, rows[i].Has("manu.id"), `The return value should have manu.id field")`)
			assert.True(t, rows[i].Has("manu.name"), `The return value should have manu.name field")`)
			assert.True(t, rows[i].Has("manu.intro"), `The return value should have not manu.intro field")`)
		}
	}

	rows = car.Reset().With("manu", func(model model.Basic) {
		model.QueryBuilder().
			Select("id", "name", "type")
	}).MustGet()
	assert.Equal(t, 6, len(rows), `The return value should be 1")`)
	if len(rows) == 6 {
		for i := range rows {
			assert.Equal(t, rows[i].Get("manu_id"), rows[i].Get("manu.id"), `The return value should be 1")`)
			assert.True(t, rows[i].Has("manu.id"), `The return value should have manu.id field")`)
			assert.True(t, rows[i].Has("manu.name"), `The return value should have manu.name field")`)
			assert.False(t, rows[i].Has("manu.intro"), `The return value should have not manu.intro field")`)
		}
	}
}

func TestModelWithHasMany(t *testing.T) {
	TestFactoryMigrate(t)
	manu := model.MakeUsing(modelTestMaker, "models/manu")
	rows := manu.With("cars").MustGet()
	assert.Equal(t, 3, len(rows), `The return value should be 3")`)
	if len(rows) == 3 {
		for i := range rows {
			cars, has := rows[i].Get("cars").([]xun.R)
			assert.True(t, has, `The return value should have cars")`)
			if has {
				assert.True(t, len(cars) > 0, `The return value should have cars")`)
				for _, car := range cars {
					assert.Equal(t, rows[i].Get("id"), car.Get("manu_id"), `The return value should be 1")`)
					assert.True(t, car.Has("id"), `The return car should have id field")`)
					assert.True(t, car.Has("name"), `The return car should have name field")`)
					assert.True(t, car.Has("created_at"), `The return car should have created_at field")`)
				}
			}
		}
	}

	rows = manu.Reset().With("cars", func(model model.Basic) {
		model.QueryBuilder().
			Select("manu_id", "name").
			Where("name", "like", "%V%").
			OrWhere("name", "like", "%M%").
			OrWhere("name", "like", "%T%").
			Take(100)
	}).MustGet()
	assert.Equal(t, 3, len(rows), `The return value should be 3")`)
	if len(rows) == 3 {
		for i := range rows {
			cars, has := rows[i].Get("cars").([]xun.R)
			assert.True(t, has, `The return value should have cars")`)
			if has {
				for _, car := range cars {
					assert.Equal(t, rows[i].Get("id"), car.Get("manu_id"), `The return value should be 1")`)
					assert.True(t, car.Has("name"), `The return car should have name field")`)
					assert.False(t, car.Has("created_at"), `The return car should not have created_at field")`)
				}
			}
		}
	}
}

func TestModelWithHasOneThrough(t *testing.T) {
	TestFactoryMigrate(t)
	manu := model.MakeUsing(modelTestMaker, "models/manu")
	rows := manu.With("user").Select("name", "type").MustGet()
	assert.Equal(t, 3, len(rows), `The return value should be 3")`)
	if len(rows) == 3 {
		assert.Equal(t, int64(1), rows[0].Get("user.user_id"), `The return value should be 1")`)
		assert.Equal(t, nil, rows[1].Get("user.user_id"), `The return value should be 1")`)
		assert.Equal(t, int64(1), rows[2].Get("user.user_id"), `The return value should be 1")`)
	}
}

func TestModelWithHasManyThrough(t *testing.T) {
	TestFactoryMigrate(t)
	manu := model.MakeUsing(modelTestMaker, "models/manu")
	rows := manu.With("users").Select("name", "type").MustGet()
	assert.Equal(t, 3, len(rows), `The return value should be 3")`)
	if len(rows) == 3 {
		users1 := rows[0].Get("users").([]xun.R)
		users2 := rows[1].Get("users").([]xun.R)
		users3 := rows[2].Get("users").([]xun.R)
		assert.Equal(t, int64(1), users1[0].Get("id"), `The return value should be 1")`)
		assert.Equal(t, 0, len(users2), `The return value should be 1")`)
		assert.Equal(t, int64(1), users3[0].Get("id"), `The return value should be 1")`)
	}
}
