package xun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustToTime(t *testing.T) {
	assert.Equal(t, "2020-12-31T23:22:14", MakeTime("2020-12-31 23:22:14").MustToTime().Format("2006-01-02T15:04:05"), "the value should be 2020-12-31T23:22:14")
	assert.Equal(t, "2021-01-01T07:22:14", MakeTime(1609456934000).MustToTime().Format("2006-01-02T15:04:05"), "the value should be 2021-01-01T07:22:14")
	assert.Equal(t, "2020-12-31", MakeTime("2020-12-31").MustToTime().Format("2006-01-02"), "the value should be 2006-01-02")
	assert.Equal(t, "23:22:14", MakeTime("23:22:14").MustToTime().Format("15:04:05"), "the value should be 23:22:14")
}

func TestMakeR(t *testing.T) {
	mapstr := map[string]string{"key": "value"}
	assert.Equal(t, "value", MakeR(mapstr)["key"], `r["key"] should be "value"`)

	mapint := map[string]int{"key": 99}
	assert.Equal(t, 99, MakeR(mapint)["key"], `r["key"] should be 99`)

	mapfloat32 := map[string]float32{"key": 99.00}
	assert.Equal(t, float32(99.00), MakeR(mapfloat32)["key"], `r["key"] should be 99.00`)

	mapfloat64 := map[string]float64{"key": 99.00}
	assert.Equal(t, float64(99.00), MakeR(mapfloat64)["key"], `r["key"] should be 99.00`)

	mapany := map[string]interface{}{"key": 99.00, "name": "hello"}
	assert.Equal(t, interface{}(99.00), MakeR(mapany)["key"], `r["key"] should be 99.00`)
	assert.Equal(t, interface{}("hello"), MakeR(mapany)["name"], `r["name"] should be "hello"`)

	mapAnyS := map[string]interface{}{
		"key":   99.00,
		"name":  "hello",
		"items": []interface{}{"item1", "item2"},
		"nested": map[string]interface{}{
			"str": "hello nested",
			"mapstr": map[string]interface{}{
				"key1": []interface{}{"s1", "s2"},
				"key2": "hello nested mapstr",
			},
		},
	}

	r := MakeR(mapAnyS)
	assert.Equal(t, interface{}(99.00), r.Get("key"), `r["key"] should be 99.00`)
	assert.Equal(t, "hello", r.Get("name"), `r["name"] should be "hello"`)
	assert.Equal(t, []interface{}{"item1", "item2"}, r.Get("items"), `r["items"] should be []interface{}{"item1", "item2"}`)
	assert.Equal(t, []interface{}{"s1", "s2"}, r.Get("nested.mapstr.key1"), `r["nested.mapstr.key1"] should be []interface{}{"s1", "s2"}`)
	assert.Equal(t, "hello nested mapstr", r.Get("nested.mapstr.key2"), `r["nested.mapstr.key2"] should be "hello nested mapstr"`)
}

func TestMakeRSlice(t *testing.T) {
	type User struct {
		Email string `json:"email"`
		Vote  int    `json:"vote"`
		Items []struct {
			ID   int
			SN   string
			Name string
		} `json:"items"`
	}

	users := []User{
		{Email: "King@example.com", Vote: 4},
		{Email: "Max@example.com", Vote: 5},
		{Email: "Jim@example.com", Vote: 6, Items: []struct {
			ID   int
			SN   string
			Name string
		}{
			{ID: 1, SN: "101", Name: "clothes"},
			{ID: 2, SN: "202", Name: "shoes"},
		}},
	}

	rows := MakeRows(users)
	assert.Equal(t, 3, len(rows), "The return rows should have 3 items")
	if len(rows) == 3 {
		assert.Equal(t, "King@example.com", rows[0].Get("email"), `rows[0]["email"] should be "King@example.com"`)
		assert.Equal(t, []R{}, rows[0].Get("items"), `rows[0]["items"] should be []R{}`)
		assert.Equal(t, interface{}(4), rows[0].Get("vote"), `rows[0]["vote"] should be 4`)

		assert.Equal(t, "Jim@example.com", rows[2].Get("email"), `rows[2]["email"] should be "Jim@example.com"`)
		assert.Equal(t, R{"id": 1, "name": "clothes", "sn": "101"}, rows[2].Get("items").([]R)[0], `rows[0]["items"] should be []R{}`)
		assert.Equal(t, interface{}(6), rows[2].Get("vote"), `rows[0]["vote"] should be 6`)
	}

	rows = MakeRows(User{Email: "King@example.com", Vote: 4})
	assert.Equal(t, 1, len(rows), "The return rows should have 1 item")
	if len(rows) == 3 {
		assert.Equal(t, "King@example.com", rows[0].Get("email"), `rows[0]["email"] should be "King@example.com"`)
		assert.Equal(t, []R{}, rows[0].Get("items"), `rows[0]["items"] should be []R{}`)
		assert.Equal(t, interface{}(4), rows[0].Get("vote"), `rows[0]["vote"] should be 4`)
	}
}
