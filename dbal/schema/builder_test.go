package schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"
)

var builder Schema

func init() {
	unit.SetLogger()
}

func TestCreate(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.Create("table_test_builder", func(table Blueprint) {
		table.ID("id").Primary()
		table.UnsignedBigInteger("counter").Index()
		table.BigInteger("latest").Index()
		table.String("name", 20).Index()
		table.String("unionid", 128).Unique()
		table.CreateUnique("name_latest", "name", "latest")
		table.CreateIndex("name_counter", "name", "counter")
	})
	assert.True(t, builder.HasTable("table_test_builder"), "should return true")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestGet(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table, err := builder.Get("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	assert.True(t, table != nil, "the return table should be BluePrint")
	if table == nil {
		return
	}
	checkTable(t, table)
}

func TestDrop(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.Drop("table_test_builder")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestDropIfExistsTableNotExists(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.DropIfExists("table_not_exists")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestDropIfExistsTableExists(t *testing.T) {
	defer unit.Catch()
	TestCreate(t)
	builder := getTestBuilder()
	err := builder.DropIfExists("table_test_builder")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestRename(t *testing.T) {
	defer unit.Catch()
	TestCreate(t)
	builder := getTestBuilder()
	err := builder.Rename("table_test_builder", "table_test_builder_re")
	assert.True(t, builder.HasTable("table_test_builder_re"), "should return true")
	assert.Equal(t, nil, err, "the return error should be nil")
	builder.Drop("table_test_builder_re")
}

func TestAlter(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropIfExists("table_test_builder")
	TestCreate(t)
	err := builder.Alter("table_test_builder", func(table Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid_unique")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
		table.CreateIndex("nickname_index", "nickname")
		table.RenameIndex("latest_index", "re_latest_index")
	})
	assert.Equal(t, nil, err, "the return error should be nil")
	if err != nil {
		return
	}
	// cheking the schema structure
	table, err := builder.Get("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	checkTableAlter(t, table)
	builder.Drop("table_test_builder")
}

func TestMustCreate(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropIfExists("table_test_builder")
	table := builder.MustCreate("table_test_builder", func(table Blueprint) {
		table.ID("id").Primary()
		table.UnsignedBigInteger("counter").Index()
		table.BigInteger("latest").Index()
		table.String("name", 20).Index()
		table.String("unionid", 128).Unique()
		table.CreateUnique("name_latest", "name", "latest")
		table.CreateIndex("name_counter", "name", "counter")
	})
	assert.True(t, builder.HasTable("table_test_builder"), "should return true")
	assert.Equal(t, "table_test_builder", table.GetName(), "the table name should be table_test_builder")
}

func TestMustGet(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table := builder.MustGet("table_test_builder")
	assert.True(t, table != nil, "the return table should be BluePrint")
	if table == nil {
		return
	}
	checkTable(t, table)
}

func TestMustDrop(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.MustDrop("table_test_builder")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
}

func TestMustDropIfExistsTableNotExists(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.MustDropIfExists("table_not_exists")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
}

func TestMustDropIfExistsTableExists(t *testing.T) {
	defer unit.Catch()
	TestMustCreate(t)
	builder := getTestBuilder()
	builder.MustDropIfExists("table_test_builder")
	assert.False(t, builder.HasTable("table_test_builder"), "should return false")
}

func TestMustRename(t *testing.T) {
	defer unit.Catch()
	TestCreate(t)
	builder := getTestBuilder()
	table := builder.MustRename("table_test_builder", "table_test_builder_re")
	assert.True(t, builder.HasTable("table_test_builder_re"), "should return true")
	assert.Equal(t, "table_test_builder_re", table.GetName(), "the table name should be table_test_builder_re")
	builder.Drop("table_test_builder_re")
}

func TestMustAlter(t *testing.T) {
	defer unit.Catch()
	TestCreate(t)
	builder := getTestBuilder()
	table := builder.MustAlter("table_test_builder", func(table Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid_unique")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
		table.CreateIndex("nickname_index", "nickname")
		table.RenameIndex("latest_index", "re_latest_index")
	})
	assert.True(t, builder.HasTable("table_test_builder"), "should return true")
	assert.Equal(t, "table_test_builder", table.GetName(), "the table name should be table_test_builder")

	// cheking the schema structure
	table, err := builder.Get("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	checkTableAlter(t, table)
	builder.Drop("table_test_builder")
}

func getTestBuilder() Schema {
	defer unit.Catch()
	if builder != nil {
		return builder
	}
	driver := os.Getenv("XUN_UNIT_DSN")
	dsn := unit.DSN(driver)
	builder = New(driver, dsn)
	return builder
}

func checkTableAlter(t *testing.T, table Blueprint) {
	DefaultBigIntPrecision := 20
	if unit.Is("postgres") {
		DefaultBigIntPrecision = 64
	}

	// checking the table schema structure
	assert.True(t, nil != table.GetColumn("id"), "the column id should be created")
	if table.GetColumn("id") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("id").Type, "the id type should be bigInteger")
		assert.Equal(t, "AutoIncrement", utils.StringVal(table.GetColumn("id").Extra), "the id extra should be AutoIncrement")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("id").Precision), "the id precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("id").IsUnsigned, "the id IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("counter"), "the column counter should be created")
	if table.GetColumn("counter") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("counter").Type, "the counter type should be bigInteger")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("counter").Precision), "the counter precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("counter").IsUnsigned, "the counter IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("latest"), "the column latest should be created")
	if table.GetColumn("latest") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("latest").Type, "the latest type should be bigInteger")
		if unit.Is("postgres") {
			assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("latest").Precision), "the latest precision should be 64")
		}
		if unit.Not("postgres") {
			assert.Equal(t, 19, utils.IntVal(table.GetColumn("latest").Precision), "the latest precision should be 19")
			assert.Equal(t, false, table.GetColumn("latest").IsUnsigned, "the latest IsUnsigned should be false")
		}
	}
	assert.True(t, nil != table.GetColumn("nickname"), "the column nickname should be created")
	if table.GetColumn("nickname") != nil {
		assert.Equal(t, "string", table.GetColumn("nickname").Type, "the name type should be string")
		assert.Equal(t, 50, utils.IntVal(table.GetColumn("nickname").Length), "the nickname length should be 50")
	}
	assert.True(t, nil != table.GetColumn("uid"), "the column uid should be created")
	if table.GetColumn("uid") != nil {
		assert.Equal(t, "string", table.GetColumn("uid").Type, "the unionid type should be string")
		if unit.Not("sqlite3") {
			assert.Equal(t, 200, utils.IntVal(table.GetColumn("uid").Length), "the unionid length should be 200")
		}
	}

	// checking the table indexes
	assert.True(t, nil != table.GetPrimary(), "the index PRIMARY should be created")
	if table.GetPrimary() != nil {
		assert.Equal(t, "id", table.GetPrimary().Columns[0].Name, "the column of PRIMARY key should be id")
	}

	assert.Equal(t, 1, len(table.GetIndex("counter_index").Columns), "the counter_index  should has 1 column")
	assert.Equal(t, "counter", table.GetIndex("counter_index").Columns[0].Name, "the column of counter_index key should be counter")
	assert.Equal(t, "index", table.GetIndex("counter_index").Type, "the counter_index key type should be index")

	if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(table.GetIndex("re_latest_index").Columns), "the re_latest_index should has 1 column")
		assert.Equal(t, "latest", table.GetIndex("re_latest_index").Columns[0].Name, "the column of re_latest_index key should be latest")
		assert.Equal(t, "index", table.GetIndex("re_latest_index").Type, "the re_latest_index key type should be index")
	}

	assert.Equal(t, 1, len(table.GetIndex("nickname_index").Columns), "the nickname_index  should has 1 column")
	assert.Equal(t, "nickname", table.GetIndex("nickname_index").Columns[0].Name, "the column of nickname_index key should be name")
	assert.Equal(t, "index", table.GetIndex("nickname_index").Type, "the nickname_index key type should be index")

	assert.Equal(t, 1, len(table.GetIndex("uid_unique").Columns), "the uid_unique should has 1 column")
	assert.Equal(t, "uid", table.GetIndex("uid_unique").Columns[0].Name, "the column of uid_unique key should be unionid")
	assert.Equal(t, "unique", table.GetIndex("uid_unique").Type, "the uid_unique key type should be unique")

	nicknameIndex := table.GetIndex("nickname_index")
	if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(nicknameIndex.Columns), "the index nickname_index should has one column")
		assert.Equal(t, "index", nicknameIndex.Type, "the nickname_index key type should be unique")
		if len(nicknameIndex.Columns) == 1 {
			assert.Equal(t, "nickname", nicknameIndex.Columns[0].Name, "the second column of the index nickname_index should be latest")
		}
	}

	reLatestIndex := table.GetIndex("re_latest_index")
	if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(reLatestIndex.Columns), "the index re_latest_index  should has one column")
		assert.Equal(t, "index", reLatestIndex.Type, "the re_latest_index key type should be unique")
		if len(reLatestIndex.Columns) == 1 {
			assert.Equal(t, "latest", reLatestIndex.Columns[0].Name, "the second column of the index re_latest_index should be latest")
		}
	}

	nameLatest := table.GetIndex("name_latest")
	if unit.Is("postgres") {
		assert.Equal(t, 0, len(nameLatest.Columns), "the index name_latest should has none")
	} else if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(nameLatest.Columns), "the index name_latest  should has one column")
		assert.Equal(t, "unique", nameLatest.Type, "the name_latest key type should be unique")
		if len(nameLatest.Columns) == 1 {
			assert.Equal(t, "latest", nameLatest.Columns[0].Name, "the second column of the index name_latest should be latest")
		}
	}

	nameCounter := table.GetIndex("name_counter")
	if unit.Is("postgres") {
		assert.Equal(t, 0, len(nameCounter.Columns), "the index name_counter should has none")
	} else if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(nameCounter.Columns), "the index name_counter should has one column")
		assert.Equal(t, "index", nameCounter.Type, "the name_counter key type should be unique")
		if len(nameCounter.Columns) == 2 {
			assert.Equal(t, "counter", nameCounter.Columns[0].Name, "the second column of the index name_counter should be counter")
		}
	}
}

func checkTable(t *testing.T, table Blueprint) {

	DefaultBigIntPrecision := 20
	if unit.Is("postgres") {
		DefaultBigIntPrecision = 64
	}

	// checking the table schema structure
	assert.True(t, nil != table.GetColumn("id"), "the column id should be created")
	if table.GetColumn("id") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("id").Type, "the id type should be bigInteger")
		assert.Equal(t, "AutoIncrement", utils.StringVal(table.GetColumn("id").Extra), "the id extra should be AutoIncrement")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("id").Precision), "the id precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("id").IsUnsigned, "the id IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("counter"), "the column counter should be created")
	if table.GetColumn("counter") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("counter").Type, "the counter type should be bigInteger")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("counter").Precision), "the counter precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("counter").IsUnsigned, "the counter IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("latest"), "the column latest should be created")
	if table.GetColumn("latest") != nil {
		assert.Equal(t, "bigInteger", table.GetColumn("latest").Type, "the latest type should be bigInteger")
		if unit.Is("postgres") {
			assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("latest").Precision), "the latest precision should be 19")
		}
		if unit.Not("postgres") {
			assert.Equal(t, 19, utils.IntVal(table.GetColumn("latest").Precision), "the latest precision should be 19")
			assert.Equal(t, false, table.GetColumn("latest").IsUnsigned, "the latest IsUnsigned should be false")
		}
	}
	assert.True(t, nil != table.GetColumn("name"), "the column name should be created")
	if table.GetColumn("name") != nil {
		assert.Equal(t, "string", table.GetColumn("name").Type, "the name type should be string")
		assert.Equal(t, 20, utils.IntVal(table.GetColumn("name").Length), "the name length should be 20")
	}
	assert.True(t, nil != table.GetColumn("unionid"), "the column unionid should be created")
	if table.GetColumn("unionid") != nil {
		assert.Equal(t, "string", table.GetColumn("unionid").Type, "the unionid type should be string")
		assert.Equal(t, 128, utils.IntVal(table.GetColumn("unionid").Length), "the unionid length should be 128")
	}

	// checking the table indexes
	assert.True(t, nil != table.GetPrimary(), "the index PRIMARY should be created")
	if table.GetPrimary() != nil {
		assert.Equal(t, "id", table.GetPrimary().Columns[0].Name, "the column of PRIMARY key should be id")
	}

	assert.Equal(t, 1, len(table.GetIndex("counter_index").Columns), "the counter_index  should has 1 column")
	assert.Equal(t, "counter", table.GetIndex("counter_index").Columns[0].Name, "the column of counter_index key should be counter")
	assert.Equal(t, "index", table.GetIndex("counter_index").Type, "the counter_index key type should be index")

	assert.Equal(t, 1, len(table.GetIndex("latest_index").Columns), "the latest_index  should has 1 column")
	assert.Equal(t, "latest", table.GetIndex("latest_index").Columns[0].Name, "the column of latest_index key should be latest")
	assert.Equal(t, "index", table.GetIndex("latest_index").Type, "the latest_index key type should be index")

	assert.Equal(t, 1, len(table.GetIndex("name_index").Columns), "the name_index  should has 1 column")
	assert.Equal(t, "name", table.GetIndex("name_index").Columns[0].Name, "the column of name_index key should be name")
	assert.Equal(t, "index", table.GetIndex("name_index").Type, "the name_index key type should be index")

	assert.Equal(t, 1, len(table.GetIndex("unionid_unique").Columns), "the unionid_unique should has 1 column")
	assert.Equal(t, "unionid", table.GetIndex("unionid_unique").Columns[0].Name, "the column of unionid_unique key should be unionid")
	assert.Equal(t, "unique", table.GetIndex("unionid_unique").Type, "the unionid_unique key type should be unique")

	nameLatest := table.GetIndex("name_latest")
	assert.Equal(t, 2, len(nameLatest.Columns), "the index name_latest  should has two columns")
	assert.Equal(t, "unique", nameLatest.Type, "the name_latest key type should be unique")
	if len(nameLatest.Columns) == 2 {
		assert.Equal(t, "name", nameLatest.Columns[0].Name, "the first column of the index name_latest  should be name")
		assert.Equal(t, "latest", nameLatest.Columns[1].Name, "the second column of the index name_latest should be latest")
	}

	nameCounter := table.GetIndex("name_counter")
	assert.Equal(t, 2, len(nameCounter.Columns), "the index name_counter should has two columns")
	assert.Equal(t, "index", nameCounter.Type, "the name_counter key type should be unique")
	if len(nameCounter.Columns) == 2 {
		assert.Equal(t, "name", nameCounter.Columns[0].Name, "the first column of the index name_counter should be name")
		assert.Equal(t, "counter", nameCounter.Columns[1].Name, "the second column of the index name_counter should be counter")
	}
}
