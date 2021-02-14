package schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

var builder schema.Schema

func init() {
	unit.SetLogger()
}

func getTestBuilder() schema.Schema {
	defer unit.Catch()
	if builder != nil {
		return builder
	}
	driver := os.Getenv("XUN_UNIT_DSN")
	dsn := unit.DSN(driver)
	builder = schema.New(driver, dsn)
	return builder
}

func TestCreate(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.Create("table_test_builder", func(table schema.Blueprint) {
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

	// checking the table schema sturcture
	assert.Equal(t, "bigInteger", table.GetColumn("id").Type, "the id type should be bigInteger")
	assert.Equal(t, "bigInteger", table.GetColumn("counter").Type, "the counter type should be bigInteger")
	assert.Equal(t, "bigInteger", table.GetColumn("latest").Type, "the latest type should be bigInteger")
	assert.Equal(t, "string", table.GetColumn("name").Type, "the name type should be string")
	assert.Equal(t, 20, *table.GetColumn("name").Length, "the name length should be 20")
	assert.Equal(t, "string", table.GetColumn("unionid").Type, "the unionid type should be string")
	assert.Equal(t, 128, *table.GetColumn("unionid").Length, "the unionid length should be 128")

	// checking the table indexes
	assert.Equal(t, "id", table.GetIndex("PRIMARY").Columns[0].Name, "the column of PRIMARY key should be id")
	assert.Equal(t, "primary", table.GetIndex("PRIMARY").Type, "the PRIMARY key type should be primary")

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
	TestCreate(t)
	builder := getTestBuilder()
	err := builder.Alter("table_test_builder", func(table schema.Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
	})
	assert.Equal(t, nil, err, "the return error should be nil")
	builder.Drop("table_test_builder")
}

func TestMustCreate(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table := builder.MustCreate("table_test_builder", func(table schema.Blueprint) {
		table.String("name", 20).Unique()
		table.String("unionid", 128).Unique()
	})
	assert.True(t, builder.HasTable("table_test_builder"), "should return true")
	assert.Equal(t, "table_test_builder", table.GetName(), "the table name should be table_test_builder")
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
	table := builder.MustAlter("table_test_builder", func(table schema.Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
	})
	assert.True(t, builder.HasTable("table_test_builder"), "should return true")
	assert.Equal(t, "table_test_builder", table.GetName(), "the table name should be table_test_builder")
	builder.Drop("table_test_builder")
}
