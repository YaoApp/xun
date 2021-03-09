package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/yaoapp/xun/grammar/mysql"    // Load the MySQL Grammar
	_ "github.com/yaoapp/xun/grammar/postgres" // Load the Postgres Grammar
	_ "github.com/yaoapp/xun/grammar/sqlite3"  // Load the SQLite3 Grammar
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

type columnFunc func(table Blueprint, name string, args ...int) *Column

func TestBlueprintSmallInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.SmallInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "integer", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3"), t, "smallInteger", nil)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.SmallInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "smallInteger", nil)
}

func TestBlueprintUnsignedSmallInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedSmallInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "integer", testCheckUnsigned)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "smallInteger", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "smallInteger", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.UnsignedSmallInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "smallInteger", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "smallInteger", testCheckUnsigned)
}

func TestBlueprintSmallIncrements(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		if unit.Is("sqlite3") {
			table.SmallIncrements("id").Primary()
		} else {
			table.SmallIncrements("id").Unique()
		}
		table.String("field1st", 60)
	})
	assert.Nil(t, err, "the create method should be return nil")

	table := builder.MustGet("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	if unit.Is("postgres") || unit.Is("sqlite3") {
		assert.Equal(t, "integer", column.Type, "the column type should be integer")
	} else {
		assert.Equal(t, "smallInteger", column.Type, "the column type should be smallInteger")
	}

	// Checking the index
	if unit.Is("sqlite3") {
		primary := table.GetPrimary()
		assert.Equal(t, "id", primary.Columns[0].Name, "the primary key should has the id column")

	} else {
		index := table.GetIndex("id_unique")
		assert.Equal(t, "unique", index.Type, "the id_unique type should be unique")
	}
}

func TestBlueprintInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Integer(name) })
	testCheckColumnsAfterCreate(true, t, "integer", nil)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.Integer(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "integer", nil)
}

func TestBlueprintUnsignedInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedInteger(name) })
	testCheckColumnsAfterCreate(true, t, "integer", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "integer", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.UnsignedInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "integer", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "integer", testCheckUnsigned)
}

func TestBlueprintIncrements(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		if unit.Is("sqlite3") {
			table.Increments("id").Primary()
		} else {
			table.Increments("id").Unique()
		}
		table.String("field1st", 60)
	})
	assert.Nil(t, err, "the create method should be return nil")

	table := builder.MustGet("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)
	assert.Equal(t, "integer", column.Type, "the column type should be integer")

	// Checking the index
	if unit.Is("sqlite3") {
		primary := table.GetPrimary()
		assert.Equal(t, "id", primary.Columns[0].Name, "the primary key should has the id column")

	} else {
		index := table.GetIndex("id_unique")
		assert.Equal(t, "unique", index.Type, "the id_unique type should be unique")
	}
}

func TestBlueprintBigInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "integer", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3"), t, "bigInteger", nil)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "bigInteger", nil)
}

func TestBlueprintUnsignedBigInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedBigInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "integer", testCheckUnsigned)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "bigInteger", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "bigInteger", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.UnsignedBigInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "bigInteger", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "bigInteger", testCheckUnsigned)
}

func TestBlueprintBigIncrements(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		if unit.Is("sqlite3") {
			table.BigIncrements("id").Primary()
		} else {
			table.BigIncrements("id").Unique()
		}
		table.String("field1st", 60)
	})
	assert.Nil(t, err, "the create method should be return nil")

	table := builder.MustGet("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	// Checking the index
	if unit.Is("sqlite3") {
		primary := table.GetPrimary()
		assert.Equal(t, "id", primary.Columns[0].Name, "the primary key should has the id column")
		assert.Equal(t, "integer", column.Type, "the column type should be integer")
	} else {
		index := table.GetIndex("id_unique")
		assert.Equal(t, "unique", index.Type, "the id_unique type should be unique")
		assert.Equal(t, "bigInteger", column.Type, "the column type should be bigInteger")
	}
}

func TestBlueprintID(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.String("field1st", 60)
	})
	assert.Nil(t, err, "the create method should be return nil")

	table := builder.MustGet("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	// Checking the index
	primary := table.GetPrimary()
	assert.Equal(t, "id", primary.Columns[0].Name, "the primary key should has the id column")
}

func TestBlueprintDecimal(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.Decimal(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "decimal", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column {
			total := 10
			places := 2
			if len(args) >= 2 {
				total = args[0] + args[1]
				places = args[1]
			}
			return table.Decimal(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "decimal", nil)
}

func TestBlueprintUnsignedDecimal(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.UnsignedDecimal(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "decimal", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "decimal", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column {
			total := 10
			places := 2
			if len(args) >= 2 {
				total = args[0] + args[1]
				places = args[1]
			}
			return table.UnsignedDecimal(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "decimal", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "decimal", testCheckUnsigned)
}

func TestBlueprinString(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) })
	testCheckColumnsAfterCreate(unit.Always, t, "string", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "string", nil)
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

func testAlterTable(executable bool, t *testing.T, create columnFunc, alter columnFunc) {
	if !executable {
		return
	}
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

func testCheckUnsigned(t *testing.T, name string, column *Column) {
	assert.True(t, column.IsUnsigned, "the column %s IsUnsigned should be true", name)
}

func testCheckAutoIncrementing(t *testing.T, name string, column *Column) {
	assert.NotNil(t, column, "the column %s should not be nil", name)
	if unit.Not("postgres") {
		assert.True(t, column.IsUnsigned, "the column %s IsUnsigned should be true", name)
	}
	assert.Equal(t, "AutoIncrement", utils.StringVal(column.Extra), "the column %s extra should be AutoIncrement", name)
}

func testCheckColumnsAfterCreate(executable bool, t *testing.T, typeName string, check func(t *testing.T, name string, column *Column)) {
	if !executable {
		return
	}
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
				check(t, name, column)
			}
		}
	}
}

func testCheckColumnsAfterAlter(executable bool, t *testing.T, typeName string, check func(t *testing.T, name string, column *Column)) {
	if !executable {
		return
	}

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
				check(t, name, column)
			}
		}
	}
}

func testCheckIndexesAfterCreate(executable bool, t *testing.T, check func(t *testing.T, name string, index *Index)) {
	if !executable {
		return
	}
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
			check(t, name, index)
		}
	}
}
