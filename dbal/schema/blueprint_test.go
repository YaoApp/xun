package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/yaoapp/xun/grammar/mysql"    // Load the MySQL Grammar
	_ "github.com/yaoapp/xun/grammar/postgres" // Load the Postgres Grammar
	_ "github.com/yaoapp/xun/grammar/sqlite3"  // Load the SQLite3 Grammar
	"github.com/yaoapp/xun/unit"
)

type columnFunc func(table Blueprint, name string, args ...int) *Column

func TestBlueprintSmallInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.SmallInteger(name)
	})
	if unit.Is("sqlite3") {
		testCheckColumnsAfterCreate(t, "integer", nil)
	} else {
		testCheckColumnsAfterCreate(t, "smallInteger", nil)
	}
	testCheckIndexesAfterCreate(t, nil)

	if unit.Not("sqlite3") {
		testAlterTable(t,
			func(table Blueprint, name string, args ...int) *Column {
				return table.String(name, args[0])
			},
			func(table Blueprint, name string, args ...int) *Column {
				return table.SmallInteger(name)
			},
		)

		testCheckColumnsAfterAlter(t, "smallInteger", nil)
	}
}

// clean the test data
func TestBlueprintClean(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
}

// The blueprint test utils

func testCreateTable(t *testing.T, create columnFunc) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		create(table, "field", 1, 1)
		create(table, "field2nd", 2, 2)
		create(table, "field3rd", 4, 4)
		create(table, "field4th", 8, 6)
		create(table, "field5th", 16, 8)
		create(table, "field6th", 32, 16)
		create(table, "field7th", 64, 32)
		create(table, "field8th", 128, 64)
		create(table, "field9th", 256, 128)
		create(table, "fieldWithIndex", 32, 2).Index()
		create(table, "fieldWithUnique", 64, 4).Unique()
		table.AddUnique("field_field2nd", "field", "field2nd")
		table.AddIndex("field2nd_field3rd", "field2nd", "field3rd")
	})
	assert.Equal(t, nil, err, "the return error should be nil")
}

func testAlterTable(t *testing.T, create columnFunc, alter columnFunc) {
	testCreateTable(t, create)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		alter(table, "field1st", 1, 1)         // Create new column
		alter(table, "field2nd", 4, 4)         // Alter field2nd column
		alter(table, "field4th", 16, 8)        // Alter field4th column
		alter(table, "fieldWithIndex", 32, 2)  // Alter fieldWithIndex column
		alter(table, "fieldWithUnique", 64, 4) // Alter fieldWithIndex column
	})
	assert.Equal(t, nil, err, "the return error should be nil")
}

func testGetTable() Blueprint {
	builder := getTestBuilder()
	return builder.MustGet("table_test_blueprint")
}

func testCheckColumnsAfterCreate(t *testing.T, typeName string, check func(name string, column *Column)) {
	table := testGetTable()
	columns := table.GetColumns()
	names := []string{
		"field", "field2nd", "field3rd", "field4th", "field5th", "field6th", "field7th", "field8th", "field9th",
		"fieldWithIndex", "fieldWithUnique",
	}
	assert.True(t, table.HasColumn(names...), "should return true")
	for name, column := range columns {
		if name != "id" {
			assert.Equal(t, typeName, column.Type, "the column type should be %s", typeName)
			if check != nil {
				check(name, column)
			}
		}
	}
}

func testCheckColumnsAfterAlter(t *testing.T, typeName string, check func(name string, column *Column)) {
	table := testGetTable()
	alterNames := []string{
		"field1st", "field2nd", "field4th", "fieldWithIndex", "fieldWithUnique",
	}
	names := []string{
		"field", "field1st", "field2nd", "field3rd", "field4th", "field5th", "field6th", "field7th", "field8th", "field9th",
		"fieldWithIndex", "fieldWithUnique",
	}
	assert.True(t, table.HasColumn(names...), "should return true")
	for _, name := range alterNames {
		column := table.GetColumn(name)
		if name != "id" {
			assert.Equal(t, typeName, column.Type, "the column type should be %s", typeName)
			if check != nil {
				check(name, column)
			}
		}
	}
}

func testCheckIndexesAfterCreate(t *testing.T, check func(name string, index *Index)) {
	table := testGetTable()
	indexes := table.GetIndexes()
	names := []string{
		"fieldWithIndex_index", "fieldWithUnique_unique",
		"field_field2nd", "field2nd_field3rd",
	}
	assert.True(t, table.HasIndex(names...), "should return true")
	assert.Equal(t, "unique", table.Index("field_field2nd").Type, "the field_field2nd type should be unique")
	assert.Equal(t, "index", table.Index("field2nd_field3rd").Type, "the field2nd_field3rd type should be index")

	fieldWithIndex := table.GetIndex("fieldWithIndex_index")
	if fieldWithIndex != nil {
		assert.Equal(t, "index", fieldWithIndex.Type, "the fieldWithIndex_index type should be index")
		assert.Equal(t, 1, len(fieldWithIndex.Columns), "the fieldWithIndex_index should has one column")
		assert.Equal(t, "fieldWithIndex", fieldWithIndex.Columns[0].Name, "the fieldWithIndex_index contains column name should be fieldWithIndex")
	}

	fieldWithUnique := table.GetIndex("fieldWithUnique_unique")
	if fieldWithUnique != nil {
		assert.Equal(t, "unique", fieldWithUnique.Type, "the fieldWithUnique_unique type should be unique")
		assert.Equal(t, 1, len(fieldWithUnique.Columns), "the fieldWithUnique_unique should has one column")
		assert.Equal(t, "fieldWithUnique", fieldWithUnique.Columns[0].Name, "the fieldWithUnique_unique contains column name should be fieldWithUnique")
	}

	fieldField2nd := table.GetIndex("field_field2nd")
	if fieldField2nd != nil {
		assert.Equal(t, "unique", fieldField2nd.Type, "the field_field2nd type should be unique")
		assert.Equal(t, 2, len(fieldField2nd.Columns), "the field_field2nd should has two columns")
	}

	field2ndField3rd := table.GetIndex("field2nd_field3rd")
	if fieldField2nd != nil {
		assert.Equal(t, "index", field2ndField3rd.Type, "the field2nd_field3rd type should be index")
		assert.Equal(t, 2, len(field2ndField3rd.Columns), "the field2nd_field3rd should has two columns")
	}

	for name, index := range indexes {
		if check != nil {
			check(name, index)
		}
	}
}
