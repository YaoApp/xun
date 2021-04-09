package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestInsertMustInsert(t *testing.T) {

	NewTableForInsertTest()
	qb := getTestBuilder()

	type User struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}
	users := []User{}

	qb.Table("table_test_insert").MustInsert(xun.R{
		"email": "kayla@example.com", "vote": 0,
	})
	err := qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 1, len(users), "The return users should be 1")
	if len(users) == 1 {
		assert.Equal(t, "kayla@example.com", users[0].Email, "The email of the first row should be kayla@example.com")
		assert.Equal(t, 0, users[0].Vote, "The vote  of the first row should be 0")
	}

	qb.Table("table_test_insert").MustInsert([]xun.R{
		{"email": "picard@example.com", "vote": 1},
		{"email": "janeway@example.com", "vote": 2},
	})
	users = nil
	err = qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 3, len(users), "The return users should be 3")
	if len(users) == 3 {
		assert.Equal(t, "kayla@example.com", users[0].Email, "The email of the first row should be kayla@example.com")
		assert.Equal(t, 0, users[0].Vote, "The vote  of the first row should be 0")
		assert.Equal(t, "picard@example.com", users[1].Email, "The email of the first row should be picard@example.com")
		assert.Equal(t, 1, users[1].Vote, "The vote  of the first row should be 1")
		assert.Equal(t, "janeway@example.com", users[2].Email, "The email of the first row should be  janeway@example.com")
		assert.Equal(t, 2, users[2].Vote, "The vote  of the first row should be 2")
	}

	user := struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}{
		Email: "Jim@example.com", Vote: 3,
	}
	qb.Table("table_test_insert").MustInsert(user)
	users = nil
	err = qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 4, len(users), "The return users should be 4")
	if len(users) == 4 {
		assert.Equal(t, "Jim@example.com", users[3].Email, "The email of the first row should be Jim@example.com")
		assert.Equal(t, 3, users[3].Vote, "The vote  of the first row should be 3")
	}

	users = []User{
		{Email: "King@example.com", Vote: 4},
		{Email: "Max@example.com", Vote: 5},
	}
	qb.Table("table_test_insert").MustInsert(users)
	users = nil
	err = qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 6, len(users), "The return users should be 4")
	if len(users) == 6 {
		assert.Equal(t, "King@example.com", users[4].Email, "The email of the first row should be King@example.com")
		assert.Equal(t, 4, users[4].Vote, "The vote  of the first row should be 4")
		assert.Equal(t, "Max@example.com", users[5].Email, "The email of the first row should be Max@example.com")
		assert.Equal(t, 5, users[5].Vote, "The vote  of the first row should be 5")
	}

}

func TestInsertMustInsertRaw(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	raw := "random()"
	if unit.DriverIs("mysql") {
		raw = "rand()"
	}

	qb.Table("table_test_insert").MustInsert([][]interface{}{
		{"picard@example.com", dbal.Raw(raw)},
		{"janeway@example.com", 2},
	}, []string{"email", "vote"})

	users := qb.Select("id", "email", "vote").OrderBy("id").MustGet()

	assert.Equal(t, 2, len(users), "The return users should be 2")
	if len(users) == 2 {
		assert.Equal(t, "picard@example.com", users[0]["email"].(string), "The email of the first row should be picard@example.com")
		assert.Equal(t, "janeway@example.com", users[1]["email"].(string), "The email of the second row should be janeway@example.com")
	}
}

func TestInsertMustInsertWithColumns(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	qb.Table("table_test_insert").MustInsert([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, []string{"email", "vote"})

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertWithColumnsStyle2(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	qb.Table("table_test_insert").MustInsert([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, "email,vote")

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertWithColumnsStyle3(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	qb.Table("table_test_insert").MustInsert([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, []string{"email,vote"})

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertWithColumnsStyle4(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	qb.Table("table_test_insert").MustInsert([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, "email", "vote")

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertOrIgnore(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()

	type User struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}
	users := []User{
		{Email: "Max@example.com", Vote: 3},
		{Email: "Max@example.com", Vote: 6},
		{Email: "King@example.com", Vote: 9},
	}

	affected := qb.Table("table_test_insert").MustInsertOrIgnore(users)
	assert.Equal(t, int64(2), affected, "The affected rows should be 2")

	users = nil
	err := qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 2, len(users), "The return users should be 2")
	if len(users) == 2 {
		assert.Equal(t, "Max@example.com", users[0].Email, "The email of the first row should be picard@example.com")
		assert.Equal(t, 3, users[0].Vote, "The vote of the first row should be 3")
		assert.Equal(t, "King@example.com", users[1].Email, "The email of the first row should be King@example.com")
		assert.Equal(t, 9, users[1].Vote, "The vote of the first row should be 9")
	}

	newQuery := New(unit.Driver(), unit.DSN())
	newQuery.DB().Close()
	_, err = newQuery.Table("table_test_insert").InsertOrIgnore(users)
	assert.Error(t, err, "the return value sholud be error")
}

func TestInsertMustInsertOrIgnoreWithColumns(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	qb.Table("table_test_insert").MustInsertOrIgnore([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, "email,vote")

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertGetID(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()

	type User struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}
	users := []User{
		{Email: "Max@example.com", Vote: 3},
		{Email: "King@example.com", Vote: 9},
	}
	qb.Table("table_test_insert").MustInsert(users)

	id := qb.Table("table_test_insert").MustInsertGetID(User{
		Email: "Jim@example.com", Vote: 7,
	})
	assert.Equal(t, int64(3), id, "The return last id should be 3")

	id = qb.Table("table_test_insert").MustInsertGetID([]User{
		{Email: "Tom@example.com", Vote: 8},
		{Email: "Gee@example.com", Vote: 10},
	})

	if unit.DriverIs("sqlite3") {
		assert.Equal(t, int64(5), id, "The return last id should be 5")
	} else {
		assert.Equal(t, int64(4), id, "The return last id should be 4")
	}

	id = qb.Table("table_test_insert").MustInsertGetID(User{
		Email: "Bee@example.com", Vote: 12,
	}, "id")
	assert.Equal(t, int64(6), id, "The return last id should be 6")

	newQuery := New(unit.Driver(), unit.DSN())
	newQuery.DB().Close()
	_, err := newQuery.Table("table_test_insert").InsertGetID(users)
	assert.Error(t, err, "the return value sholud be error")
}

func TestInsertMustInsertGetIdWithColumns(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	id := qb.Table("table_test_insert").MustInsertGetID([][]interface{}{
		{"picard@example.com", 1},
		{"janeway@example.com", 2},
	}, "id", "email,vote")

	if unit.DriverIs("sqlite3") {
		assert.Equal(t, int64(2), id, "The return last id should be 2")
	} else {
		assert.Equal(t, int64(1), id, "The return last id should be 1")
	}

	checkInsertWithColumns(t, qb)
}

func TestInsertMustInsertUsing(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	var affected int64
	var sql string
	if unit.DriverIs("postgres") {
		sql = "$1 as email, $2 as vote"
	} else {
		sql = "? as email, ? as vote"
	}
	affected = qb.Table("table_test_insert").MustInsertUsing(func(qb Query) {
		qb.SelectRaw(sql, "Bee@example.com", 2)
	}, "email,vote")

	assert.Equal(t, int64(1), affected, "The return affected should be 1")
	assert.Panics(t, func() {
		newQuery := New(unit.Driver(), unit.DSN())
		newQuery.DB().Close()
		newQuery.Table("table_test_insert").MustInsertUsing(func(qb Query) {
			qb.SelectRaw(sql, "Bee@example.com", 2)
		}, "email", "vote")
	})
}

func TestInsertMustInsertUsingWithComma(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	var affected int64
	var sql string
	if unit.DriverIs("postgres") {
		sql = "$1 as email, $2 as vote"
	} else {
		sql = "? as email, ? as vote"
	}

	affected = qb.Table("table_test_insert").MustInsertUsing(func(qb Query) {
		qb.SelectRaw(sql, "Bee@example.com", 2)
	}, "email,vote")

	assert.Equal(t, int64(1), affected, "The return affected should be 1")
}

func TestInsertMustInsertUsingWithArray(t *testing.T) {
	NewTableForInsertTest()
	qb := getTestBuilder()
	var affected int64
	var sql string
	if unit.DriverIs("postgres") {
		sql = "$1 as email, $2 as vote"
	} else {
		sql = "? as email, ? as vote"
	}

	affected = qb.Table("table_test_insert").MustInsertUsing(func(qb Query) {
		qb.SelectRaw(sql, "Bee@example.com", 2)
	}, []string{"email", "vote"})

	assert.Equal(t, int64(1), affected, "The return affected should be 1")

}

// clean the test data
func TestInsertClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_insert")
}

func NewTableForInsertTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_insert")
	builder.MustCreateTable("table_test_insert", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email").Unique()
		table.Integer("vote")
	})
}

func checkInsertWithColumns(t *testing.T, qb Query) {
	users := qb.Select("email", "vote").OrderBy("vote").MustGet()
	assert.Equal(t, 2, len(users), "The return users should be 2")
	if len(users) == 2 {
		assert.Equal(t, "picard@example.com", users[0]["email"].(string), "The email of the first row should be picard@example.com")
		assert.Equal(t, int64(1), users[0]["vote"].(int64), "The vote of the first row should be 1")
		assert.Equal(t, "janeway@example.com", users[1]["email"].(string), "The email of the second row should be janeway@example.com")
		assert.Equal(t, int64(2), users[1]["vote"].(int64), "The vote of the second row should be 2")
	}
}
