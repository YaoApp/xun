package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

func TestInsertMustInsert(t *testing.T) {

	NewTableForColumnTest()

	qb := getTestBuilder()

	users := []struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}{}

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
		assert.Equal(t, "janeway@example.com", users[2].Email, "The email of the first row should be picard@example.com")
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
		assert.Equal(t, "Jim@example.com", users[3].Email, "The email of the first row should be picard@example.com")
		assert.Equal(t, 3, users[3].Vote, "The vote  of the first row should be 3")
	}

	users = []struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
	}{
		{Email: "King@example.com", Vote: 4},
		{Email: "Max@example.com", Vote: 5},
	}
	qb.Table("table_test_insert").MustInsert(users)
	users = nil
	err = qb.DB(true).Select(&users, "SELECT email,vote FROM table_test_insert LIMIT 20")
	assert.True(t, err == nil, "The return error should be nil")
	assert.Equal(t, 6, len(users), "The return users should be 4")
	if len(users) == 6 {
		assert.Equal(t, "King@example.com", users[4].Email, "The email of the first row should be picard@example.com")
		assert.Equal(t, 4, users[4].Vote, "The vote  of the first row should be 4")
		assert.Equal(t, "Max@example.com", users[5].Email, "The email of the first row should be picard@example.com")
		assert.Equal(t, 5, users[5].Vote, "The vote  of the first row should be 5")
	}

}

// clean the test data
func TestIndexClean(t *testing.T) {
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_insert")
}

func NewTableForColumnTest() {
	defer unit.Catch()
	builder := getTestSchemaBuilder()
	builder.DropTableIfExists("table_test_insert")
	builder.MustCreateTable("table_test_insert", func(table schema.Blueprint) {
		table.ID("id")
		table.String("email").Index()
		table.Integer("vote")
	})
}
