package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestPrimaryAddPrimary(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	builder.DropTableIfExists("table_test_primary")
	builder.MustCreateTable("table_test_primary", func(table Blueprint) {
		table.BigIncrements("id")
		table.String("field1", 40)
		table.String("field2", 40)
		table.AddPrimary("id")
	})
	table := builder.MustGetTable("table_test_primary")
	primryKey := table.GetPrimary()
	CheckPrimaryKey(t, primryKey)
}

func TestPrimaryAddPrimaryFail(t *testing.T) {
	defer unit.Catch()
	if unit.DriverIs("sqlite3") {
		return
	}
	builder := getTestBuilderInstance()
	builder.DropTableIfExists("table_test_primary")
	TestPrimaryAddPrimary(t)
	err := builder.AlterTable("table_test_primary", func(table Blueprint) {
		table.Text("hello")
		table.DropPrimary()
		table.AddPrimary("id", "hello", "field1", "field2")
	})

	if unit.DriverIs("mysql") {
		assert.False(t, err == nil, "The return error should not be nil")
		table := builder.MustGetTable("table_test_primary")
		primryKey := table.GetPrimary()
		assert.True(t, primryKey == nil, "The primary key should be nil")
	} else if unit.DriverIs("postgres") {
		assert.True(t, err == nil, "The return error should  be nil")
		table := builder.MustGetTable("table_test_primary")
		primryKey := table.GetPrimary()
		assert.Equal(t, 4, len(primryKey.Columns), "The primary key should have 4 columns")
	}
}

func TestPrimaryGetPrimary(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilderInstance()
	TestPrimaryAddPrimary(t)
	table := builder.MustGetTable("table_test_primary")
	primryKey := table.GetPrimary()
	CheckPrimaryKey(t, primryKey)
}

func TestPrimaryDropPrimary(t *testing.T) {
	defer unit.Catch()
	if unit.DriverIs("sqlite3") {
		return
	}
	builder := getTestBuilderInstance()
	TestPrimaryAddPrimary(t)
	builder.MustAlterTable("table_test_primary", func(table Blueprint) {
		table.DropPrimary()
	})

	table := builder.MustGetTable("table_test_index")
	primryKey := table.GetPrimary()
	assert.True(t, primryKey == nil, "The primary key should be nil")

	table = builder.MustGetTable("table_test_primary")
	primryKey = table.GetPrimary()
	assert.True(t, primryKey == nil, "The primary key should be nil")
}

// clean the test data
func TestPrimaryClean(t *testing.T) {
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_primary")
}

func CheckPrimaryKey(t *testing.T, primryKey *Primary) {
	assert.False(t, primryKey == nil, "The primary key should be nil")
	if primryKey == nil {
		return
	}
	columnCount := len(primryKey.Columns)
	assert.Equal(t, 1, columnCount, "The primary key should have one column")
	if columnCount == 1 {
		id := primryKey.Columns[0]
		assert.Equal(t, "id", id.Name, "The primary key should contains id column")
	}
	assert.Equal(t, "PRIMARY", primryKey.Name, "The primary key name should be pk_id")
}
