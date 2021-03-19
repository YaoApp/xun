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

func TestBlueprintTinyInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.TinyInteger(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "tinyInteger", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "smallInteger", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.TinyInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "tinyInteger", nil)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "smallInteger", nil)
}

func TestBlueprintUnsignedTinyInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedTinyInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "tinyInteger", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "smallInteger", nil)
	testCheckColumnsAfterCreate(unit.Not("postgres") && unit.Not("sqlite3"), t, "tinyInteger", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.UnsignedTinyInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "tinyInteger", testCheckUnsigned)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "smallInteger", nil)
}

func TestBlueprintTinyIncrements(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	err := builder.Create("table_test_blueprint", func(table Blueprint) {
		if unit.Is("sqlite3") {
			table.TinyIncrements("id").Primary()
		} else {
			table.TinyIncrements("id").Unique()
		}
		table.String("field1st", 60)
	})
	assert.Nil(t, err, "the create method should be return nil")

	table := builder.MustGetTable("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	if unit.Is("postgres") {
		assert.Equal(t, "smallInteger", column.Type, "the column type should be smallInteger")
	} else if unit.Is("sqlite3") {
		assert.Equal(t, "integer", column.Type, "the column type should be integer")
	} else {
		assert.Equal(t, "tinyInteger", column.Type, "the column type should be tinyInteger")
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

func TestBlueprintSmallInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.SmallInteger(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "smallInteger", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.SmallInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "smallInteger", nil)
}

func TestBlueprintUnsignedSmallInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedSmallInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("postgres") || unit.Is("sqlite3"), t, "smallInteger", nil)
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

	table := builder.MustGetTable("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	if unit.Is("sqlite3") {
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

	table := builder.MustGetTable("table_test_blueprint")
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
	testCheckColumnsAfterCreate(unit.Always, t, "bigInteger", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "bigInteger", nil)
}

func TestBlueprintUnsignedBigInteger(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.UnsignedBigInteger(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "bigInteger", testCheckUnsigned)
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

	table := builder.MustGetTable("table_test_blueprint")
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

	table := builder.MustGetTable("table_test_blueprint")
	column := table.GetColumn("id")
	testCheckAutoIncrementing(t, "id", column)

	// Checking the index
	primary := table.GetPrimary()
	assert.Equal(t, "id", primary.Columns[0].Name, "the primary key should has the id column")
}

func TestBlueprintForeignID(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.ForeignID(name) })
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "bigInteger", testCheckUnsigned)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "bigInteger", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "bigInteger", testCheckUnsigned)
	testCheckIndexesAfterCreate(true, t, nil)

	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
		func(table Blueprint, name string, args ...int) *Column { return table.ForeignID(name) },
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "bigInteger", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "bigInteger", testCheckUnsigned)
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

func TestBlueprintFloat(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.Float(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "float", nil)
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
			return table.Float(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "float", nil)
}

func TestBlueprintUnsignedFloat(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.UnsignedFloat(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "float", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "float", testCheckUnsigned)
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
			return table.UnsignedFloat(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "float", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "float", testCheckUnsigned)
}

func TestBlueprintDouble(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.Double(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "double", nil)
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
			return table.Double(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "double", nil)
}

func TestBlueprintUnsignedDouble(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		total := 10
		places := 2
		if len(args) >= 2 {
			total = args[0] + args[1]
			places = args[1]
		}
		return table.UnsignedDouble(name, total, places)
	})
	testCheckColumnsAfterCreate(true, t, "double", nil)
	testCheckColumnsAfterCreate(unit.Not("sqlite3") && unit.Not("postgres"), t, "double", testCheckUnsigned)
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
			return table.UnsignedDouble(name, total, places)
		},
	)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "double", nil)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "double", testCheckUnsigned)
}

func TestBlueprintString(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) })
	testCheckColumnsAfterCreate(unit.Always, t, "string", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, args[0]) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "string", nil)
}

func TestBlueprintChar(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Char(name, args[0]) })
	testCheckColumnsAfterCreate(unit.Always, t, "char", nil)
	testCheckIndexesAfterCreate(true, t, nil)
	testAlterTable(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.Char(name, args[0]) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "char", nil)
}

func TestBlueprintText(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Text(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "text", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.Text(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "text", nil)
}

func TestBlueprintMediumText(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.MediumText(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "mediumText", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "text", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.MediumText(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "mediumText", nil)
}

func TestBlueprintLongText(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.LongText(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "longText", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "text", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.BigInteger(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.LongText(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "longText", nil)
}

func TestBlueprintBinary(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Binary(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "binary", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.Text(name) },
		func(table Blueprint, name string, args ...int) *Column { return table.Binary(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "binary", nil)
}

func TestBlueprintDate(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Date(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "date", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.Date(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "date", nil)
}

func TestBlueprintDateTime(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.DateTime(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "dateTime", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestamp", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.DateTime(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "dateTime", nil)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestamp", nil)
}

func TestBlueprintDateTimeWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.DateTime(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Not("postgres") && unit.Not("sqlite3"), t, "dateTime", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestamp", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "dateTime", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.DateTime(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "dateTime", testCheckDateTimePrecision6)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestamp", testCheckDateTimePrecision6)
}

func TestBlueprintDateTimeTz(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.DateTimeTz(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "dateTime", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestampTz", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.DateTimeTz(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "dateTime", nil)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestampTz", nil)
}

func TestBlueprintDateTimeTzWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.DateTimeTz(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Not("postgres") && unit.Not("sqlite3"), t, "dateTime", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestampTz", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "dateTime", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.DateTimeTz(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "dateTime", testCheckDateTimePrecision6)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestampTz", testCheckDateTimePrecision6)
}

func TestBlueprintTime(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Time(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "time", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.Time(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "time", nil)
}

func TestBlueprintTimeWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.Time(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Not("sqlite3"), t, "time", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "time", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.Time(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "time", testCheckDateTimePrecision6)
}

func TestBlueprintTimeTz(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.TimeTz(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "time", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timeTz", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.TimeTz(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "time", nil)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timeTz", nil)
}

func TestBlueprintTimeTzWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.TimeTz(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Not("postgres") && unit.Not("sqlite3"), t, "time", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timeTz", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "time", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.TimeTz(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "time", testCheckDateTimePrecision6)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timeTz", testCheckDateTimePrecision6)
}

func TestBlueprintTimestamp(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Timestamp(name) })
	testCheckColumnsAfterCreate(unit.Always, t, "timestamp", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.Timestamp(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "timestamp", nil)
}

func TestBlueprintTimestampWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.Timestamp(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Always, t, "timestamp", testCheckDateTimePrecision6)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.Timestamp(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "timestamp", testCheckDateTimePrecision6)
}

func TestBlueprintTimestampTz(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.TimestampTz(name) })
	testCheckColumnsAfterCreate(unit.Not("postgres"), t, "timestamp", nil)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestampTz", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column { return table.TimestampTz(name) },
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "timestamp", nil)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestampTz", nil)
}

func TestBlueprintTimestampTzWithP(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.TimestampTz(name).SetDateTimePrecision(6)
	})
	testCheckColumnsAfterCreate(unit.Not("postgres") && unit.Not("sqlite3"), t, "timestamp", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("postgres"), t, "timestampTz", testCheckDateTimePrecision6)
	testCheckColumnsAfterCreate(unit.Is("sqlite3"), t, "timestamp", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.TimestampTz(name).SetDateTimePrecision(6)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.Not("postgres"), t, "timestamp", testCheckDateTimePrecision6)
	testCheckColumnsAfterAlter(unit.Is("postgres"), t, "timestampTz", testCheckDateTimePrecision6)
}

func TestBlueprintBoolean(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column { return table.Boolean(name) })
	testCheckColumnsAfterCreate(unit.DriverNot("mysql"), t, "boolean", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("mysql"), t, "tinyInteger", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.Boolean(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3") && unit.DriverNot("mysql"), t, "boolean", nil)
	testCheckColumnsAfterAlter(unit.DriverIs("mysql"), t, "tinyInteger", nil)
}

func TestBlueprintEnum(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.Enum(name, []string{"O1", "O2", "O3"})
	})
	testCheckColumnsAfterCreate(unit.Always, t, "enum", testCheckOptionO1O2O3)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.Enum(name, []string{"O1", "O2", "O3"})
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "enum", testCheckOptionO1O2O3)
}

func TestBlueprintJSON(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.JSON(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "json", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "text", nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.JSON(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "json", nil)
}

func TestBlueprintJSONB(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.JSONB(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "jsonb", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "text", nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 128) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.JSONB(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "jsonb", nil)
}

func TestBlueprintUUID(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.UUID(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "uuid", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "string", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.Text(name) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.UUID(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "uuid", nil)
}

func TestBlueprintIPAddress(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.IPAddress(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "ipAddress", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "integer", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.IPAddress(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "ipAddress", nil)
}

func TestBlueprintMACAddress(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.MACAddress(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "macAddress", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "bigInteger", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 48) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.MACAddress(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "macAddress", nil)
}

func TestBlueprintYear(t *testing.T) {
	testCreateTable(t, func(table Blueprint, name string, args ...int) *Column {
		return table.Year(name)
	})
	testCheckColumnsAfterCreate(unit.DriverNot("sqlite3"), t, "year", nil)
	testCheckColumnsAfterCreate(unit.DriverIs("sqlite3"), t, "smallInteger", nil)
	testCheckIndexesAfterCreate(unit.Always, t, nil)
	testAlterTableSafe(unit.Not("sqlite3"), t,
		func(table Blueprint, name string, args ...int) *Column { return table.String(name, 4) },
		func(table Blueprint, name string, args ...int) *Column {
			return table.Year(name)
		},
	)
	testCheckColumnsAfterAlter(unit.Not("sqlite3"), t, "year", nil)
}

func TestBlueprintTimestamps(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.Timestamps()
	})

	table := testGetTable()
	createdAt := table.GetColumn("created_at")
	updatedAt := table.GetColumn("updated_at")
	assert.True(t, createdAt != nil, "the column created_at should be created")
	assert.True(t, updatedAt != nil, "the column updated_at should be created")

	if createdAt != nil {
		assert.Equal(t, "timestamp", createdAt.Type, "the column created_at type should be timestamp")
		assert.True(t, createdAt.Nullable, "the column created_at nullable should be true")
	}

	if updatedAt != nil {
		assert.Equal(t, "timestamp", updatedAt.Type, "the column updated_at type should be timestamp")
		assert.True(t, updatedAt.Nullable, "the column updated_at nullable should be true")
	}
}

func TestBlueprintTimestampsWithP(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.Timestamps(6)
	})

	table := testGetTable()
	createdAt := table.GetColumn("created_at")
	updatedAt := table.GetColumn("updated_at")
	assert.True(t, createdAt != nil, "the column created_at should be created")
	assert.True(t, updatedAt != nil, "the column updated_at should be created")

	if createdAt != nil {
		assert.Equal(t, "timestamp", createdAt.Type, "the column created_at type should be timestamp")
		assert.True(t, createdAt.Nullable, "the column created_at nullable should be true")
		assert.Equal(t, 6, utils.IntVal(createdAt.DateTimePrecision), "the column created_at DateTimePrecision should be 6")
	}

	if updatedAt != nil {
		assert.Equal(t, "timestamp", updatedAt.Type, "the column updated_at type should be timestamp")
		assert.True(t, updatedAt.Nullable, "the column updated_at nullable should be true")
		assert.Equal(t, 6, utils.IntVal(updatedAt.DateTimePrecision), "the column updated_at DateTimePrecision should be 6")
	}
}

func TestBlueprintDropTimestamps(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	TestBlueprintTimestamps(t)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.DropTimestamps()
	})
	assert.True(t, err == nil, "the alter method should be return nil")
	table := testGetTable()
	assert.True(t, table.GetColumn("created_at") == nil, "the column created_at should be nil")
	assert.True(t, table.GetColumn("updated_at") == nil, "the column updated_at should be nil")
}

func TestBlueprintTimestampsTz(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.TimestampsTz()
	})

	table := testGetTable()
	createdAt := table.GetColumn("created_at")
	updatedAt := table.GetColumn("updated_at")
	assert.True(t, createdAt != nil, "the column created_at should be created")
	assert.True(t, updatedAt != nil, "the column updated_at should be created")

	if createdAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", createdAt.Type, "the column created_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", createdAt.Type, "the column created_at type should be timestamp")
		}
		assert.True(t, createdAt.Nullable, "the column created_at nullable should be true")
	}

	if updatedAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", updatedAt.Type, "the column updated_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", updatedAt.Type, "the column updated_at type should be timestamp")
		}
		assert.True(t, updatedAt.Nullable, "the column updated_at nullable should be true")
	}
}

func TestBlueprintTimestampsTzWithP(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.TimestampsTz(6)
	})

	table := testGetTable()
	createdAt := table.GetColumn("created_at")
	updatedAt := table.GetColumn("updated_at")
	assert.True(t, createdAt != nil, "the column created_at should be created")
	assert.True(t, updatedAt != nil, "the column updated_at should be created")

	if createdAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", createdAt.Type, "the column created_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", createdAt.Type, "the column created_at type should be timestamp")
		}
		assert.True(t, createdAt.Nullable, "the column created_at nullable should be true")
		assert.Equal(t, 6, utils.IntVal(createdAt.DateTimePrecision), "the column created_at DateTimePrecision should be 6")
	}

	if updatedAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", updatedAt.Type, "the column updated_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", createdAt.Type, "the column created_at type should be timestamp")
		}
		assert.True(t, updatedAt.Nullable, "the column updated_at nullable should be true")
		assert.Equal(t, 6, utils.IntVal(updatedAt.DateTimePrecision), "the column updated_at DateTimePrecision should be 6")
	}
}

func TestBlueprintDropTimestampsTz(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	TestBlueprintTimestampsTz(t)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.DropTimestampsTz()
	})
	assert.True(t, err == nil, "the alter method should be return nil")
	table := testGetTable()
	assert.True(t, table.GetColumn("created_at") == nil, "the column created_at should be nil")
	assert.True(t, table.GetColumn("updated_at") == nil, "the column updated_at should be nil")
}

func TestBlueprintSoftDeletes(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.SoftDeletes()
	})

	table := testGetTable()
	deleteAt := table.GetColumn("deleted_at")
	assert.True(t, deleteAt != nil, "the column deleted_at should be created")

	if deleteAt != nil {
		assert.Equal(t, "timestamp", deleteAt.Type, "the column deleted_at type should be timestamp")
		assert.True(t, deleteAt.Nullable, "the column deleted_at nullable should be true")
	}
}

func TestBlueprintSoftDeletesWithP(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.SoftDeletes(6)
	})

	table := testGetTable()
	deleteAt := table.GetColumn("deleted_at")
	assert.True(t, deleteAt != nil, "the column deleted_at should be created")

	if deleteAt != nil {
		assert.Equal(t, "timestamp", deleteAt.Type, "the column deleted_at type should be timestamp")
		assert.True(t, deleteAt.Nullable, "the column deleted_at nullable should be true")
	}
}

func TestBlueprintDropSoftDeletes(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	TestBlueprintSoftDeletes(t)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.DropSoftDeletes()
	})
	assert.True(t, err == nil, "the alter method should be return nil")
	table := testGetTable()
	assert.True(t, table.GetColumn("deleted_at") == nil, "the column deleted_at should be nil")
}

func TestBlueprintSoftDeletesTz(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.SoftDeletesTz()
	})

	table := testGetTable()
	deleteAt := table.GetColumn("deleted_at")
	assert.True(t, deleteAt != nil, "the column deleted_at should be created")

	if deleteAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", deleteAt.Type, "the column deleted_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", deleteAt.Type, "the column deleted_at type should be timestamp")
		}
		assert.True(t, deleteAt.Nullable, "the column deleted_at nullable should be true")
	}

}

func TestBlueprintSoftDeletesTzWithP(t *testing.T) {
	builder := getTestBuilder()
	builder.DropIfExists("table_test_blueprint")
	builder.Create("table_test_blueprint", func(table Blueprint) {
		table.ID("id")
		table.SoftDeletesTz(6)
	})

	table := testGetTable()
	deleteAt := table.GetColumn("deleted_at")
	assert.True(t, deleteAt != nil, "the column deleted_at should be created")

	if deleteAt != nil {
		if unit.DriverIs("postgres") {
			assert.Equal(t, "timestampTz", deleteAt.Type, "the column deleted_at type should be timestampTz")
		} else {
			assert.Equal(t, "timestamp", deleteAt.Type, "the column deleted_at type should be timestamp")
		}
		assert.True(t, deleteAt.Nullable, "the column deleted_at nullable should be true")
		assert.Equal(t, 6, utils.IntVal(deleteAt.DateTimePrecision), "the column deleted_at DateTimePrecision should be 6")
	}
}

func TestBlueprintDropSoftDeletesTz(t *testing.T) {
	if unit.DriverIs("sqlite3") {
		return
	}
	TestBlueprintSoftDeletes(t)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.DropSoftDeletesTz()
	})
	assert.True(t, err == nil, "the alter method should be return nil")
	table := testGetTable()
	assert.True(t, table.GetColumn("deleted_at") == nil, "the column deleted_at should be nil")
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

func testAlterTableSafe(executable bool, t *testing.T, create columnFunc, alter columnFunc) {
	if !executable {
		return
	}
	testCreateTable(t, create)
	builder := getTestBuilder()
	err := builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.DropIndex("fieldWithIndex_index")
		table.DropIndex("fieldWithUnique_unique")
		table.DropIndex("field_field2nd")
		table.DropIndex("field2nd_field3rd")
		alter(table, "field1st", 1, 1)         // Create new column
		alter(table, "field2nd", 4, 4)         // Alter field2nd column
		alter(table, "field4th", 16, 8)        // Alter field4th column
		alter(table, "fieldWithIndex", 32, 2)  // Alter fieldWithIndex column
		alter(table, "fieldWithUnique", 64, 4) // Alter fieldWithIndex column
	})
	assert.Equal(t, nil, err, "the return error should be nil")

	// Add index
	err = builder.Alter("table_test_blueprint", func(table Blueprint) {
		table.AddUnique("field_field2nd", "field", "field2nd")
		table.AddIndex("field2nd_field3rd", "field2nd", "field3rd")
	})
	assert.Equal(t, nil, err, "the return error should be nil")
}

func testGetTable() Blueprint {
	builder := getTestBuilder()
	return builder.MustGetTable("table_test_blueprint")
}

func testCheckUnsigned(t *testing.T, name string, column *Column) {
	assert.True(t, column.IsUnsigned, "the column %s IsUnsigned should be true", name)
}

func testCheckOptionO1O2O3(t *testing.T, name string, column *Column) {
	assert.Len(t, column.Option, 3, "the column %s Option should has 3 items", name)
	for _, opt := range []string{"O1", "O2", "O3"} {
		assert.True(t, utils.StringHave(column.Option, opt), "the column %s Option should have %s", name, opt)
	}
}

func testCheckDateTimePrecision6(t *testing.T, name string, column *Column) {
	assert.Equal(t, 6, utils.IntVal(column.DateTimePrecision), "the column %s DateTimePrecision should be 6", name)
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
			if column == nil {
				assert.False(t, column == nil, "the column %s should not be nil", name)
				continue
			}
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
