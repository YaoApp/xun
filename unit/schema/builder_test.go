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
