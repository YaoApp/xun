package xun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustToRows(t *testing.T) {

	users := []struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
		Items []struct {
			ID   string
			SN   string
			Name string
		} `json:"items"`
	}{
		{Email: "King@example.com", Vote: 4},
		{Email: "Max@example.com", Vote: 5},
		{Email: "Jim@example.com", Vote: 6, Items: []struct {
			ID   string
			SN   string
			Name string
		}{
			{ID: "1", SN: "101", Name: "clothes"},
			{ID: "2", SN: "202", Name: "shoes"},
		}},
	}

	rows := AnyToRows(users)
	assert.Equal(t, 3, len(rows), "The return rows should have 3 items")

	if len(rows) == 3 {
		assert.Equal(t, "King@example.com", rows[0]["email"], "The first row email should be King@example.com")
		assert.Equal(t, 4, rows[0]["vote"], "The first row vote should be 4")
		assert.Equal(t, 0, len(rows[0]["items"].([]R)), "The first row items should have 0 item")

		assert.Equal(t, "Max@example.com", rows[1]["email"], "The second row email should be King@example.com")
		assert.Equal(t, 5, rows[1]["vote"], "The second row vote should be 4")
		assert.Equal(t, 0, len(rows[1]["items"].([]R)), "The second row items should have 0 item")

		assert.Equal(t, "Jim@example.com", rows[2]["email"], "The third row email should be King@example.com")
		assert.Equal(t, 6, rows[2]["vote"], "The second third vote should be 4")
		assert.Equal(t, 2, len(rows[2]["items"].([]R)), "The third row items should have 2 items")

		if len(rows[2]["items"].([]R)) == 2 {
			assert.Equal(t, "clothes", rows[2]["items"].([]R)[0]["name"], "The third row first item name should be clothes")
			assert.Equal(t, "101", rows[2]["items"].([]R)[0]["sn"], "The third row first item sn should be 101")
			assert.Equal(t, "1", rows[2]["items"].([]R)[0]["id"], "The third row first item id should be 1")

			assert.Equal(t, "clothes", rows[2]["items"].([]R)[0]["name"], "The third row first item name should be clothes")
			assert.Equal(t, "101", rows[2]["items"].([]R)[0]["sn"], "The third row first item sn should be 101")
			assert.Equal(t, "1", rows[2]["items"].([]R)[0]["id"], "The third row first item id should be 1")

			assert.Equal(t, "shoes", rows[2]["items"].([]R)[1]["name"], "The third row second item name should be shoes")
			assert.Equal(t, "202", rows[2]["items"].([]R)[1]["sn"], "The third row second item sn should be 202")
			assert.Equal(t, "2", rows[2]["items"].([]R)[1]["id"], "The third row second item id should be 2")
		}
	}

}

func TestMustToTime(t *testing.T) {
	assert.Equal(t, "2020-12-31T23:22:14", Time("2020-12-31 23:22:14").MustToTime().Format("2006-01-02T15:04:05"), "the value should be 2020-12-31T23:22:14")
	assert.Equal(t, "2021-01-01T07:22:14", Time(1609456934000).MustToTime().Format("2006-01-02T15:04:05"), "the value should be 2021-01-01T07:22:14")
	assert.Equal(t, "2020-12-31", Time("2020-12-31").MustToTime().Format("2006-01-02"), "the value should be 2006-01-02")
	assert.Equal(t, "23:22:14", Time("23:22:14").MustToTime().Format("15:04:05"), "the value should be 23:22:14")
}
