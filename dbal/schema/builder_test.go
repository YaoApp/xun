package schema

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/unit"
	"github.com/yaoapp/xun/utils"

	_ "github.com/yaoapp/xun/grammar/mysql"    // Load the MySQL Grammar
	_ "github.com/yaoapp/xun/grammar/postgres" // Load the Postgres Grammar
	_ "github.com/yaoapp/xun/grammar/sqlite3"  // Load the SQLite3 Grammar
)

var testBuilder Schema
var testBuilderInstance *Builder

func getTestBuilder() Schema {
	defer unit.Catch()
	unit.SetLogger()
	if testBuilder != nil {
		return testBuilder
	}
	driver := unit.Driver()
	dsn := unit.DSN()
	testBuilder = New(driver, dsn)
	return testBuilder
}

func getTestBuilderInstance() *Builder {
	defer unit.Catch()
	unit.SetLogger()
	if testBuilderInstance != nil {
		return testBuilderInstance
	}
	driver := unit.Driver()
	dsn := unit.DSN()
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	conn := &Connection{
		Write: db,
		WriteConfig: &dbal.Config{
			DSN:    dsn,
			Driver: driver,
			Name:   "main",
		},
		Option: &dbal.Option{},
	}

	testBuilderInstance = NewBuilder(conn)
	return testBuilderInstance
}

func TestBuilderNewFail(t *testing.T) {
	assert.Panics(t, func() {
		driver := unit.Driver()
		New(driver, "file:/root/error_dsn")
	})
}

func TestBuilderNewGrammarFail(t *testing.T) {

	driver := unit.Driver()
	dsn := unit.DSN()
	notSupportValue := "someSQL"
	shouldReturnError := fmt.Errorf("The %s driver not import", notSupportValue)
	assert.PanicsWithError(t, shouldReturnError.Error(), func() {
		db, err := sqlx.Open(driver, dsn)
		if err != nil {
			panic(err)
		}
		conn := &Connection{
			Write: db,
			WriteConfig: &dbal.Config{
				DSN:    dsn,
				Driver: driver,
				Name:   "test",
			},
			Option: &dbal.Option{},
		}
		conn.WriteConfig.Driver = notSupportValue
		NewGrammar(conn)
	})

	assert.PanicsWithError(t, "grammar setup error. (db is nil)", func() {
		db, err := sqlx.Open(driver, dsn)
		if err != nil {
			panic(err)
		}
		conn := &Connection{
			Write: db,
			WriteConfig: &dbal.Config{
				DSN:    dsn,
				Driver: driver,
				Name:   "test",
			},
			Option: &dbal.Option{},
		}
		conn.Write = nil
		NewGrammar(conn)
	})

	if unit.DriverIs("mysql") {
		assert.PanicsWithError(t, "the OnConnected event error. (sql: database is closed)", func() {
			db, err := sqlx.Open(driver, dsn)
			if err != nil {
				panic(err)
			}
			conn := &Connection{
				Write: db,
				WriteConfig: &dbal.Config{
					DSN:    dsn,
					Driver: driver,
					Name:   "test",
				},
				Option: &dbal.Option{},
			}
			conn.Write.Close()
			NewGrammar(conn)
		})
	}

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	conn := &Connection{
		Write: db,
		WriteConfig: &dbal.Config{
			DSN:    dsn,
			Driver: driver,
			Name:   "test",
		},
		Option: &dbal.Option{},
	}
	NewGrammar(conn)

}

func TestBuilderGetConnection(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	conn, err := builder.GetConnection()
	assert.Equal(t, nil, err, "the return error should be nil")
	value := ""
	err = conn.DB.Get(&value, "SELECT 'hello' ")
	assert.Equal(t, nil, err, "the return error should be nil")
	if err == nil {
		assert.Equal(t, "hello", value, "the return value should be hello")
	}
	assert.Equal(t, unit.Driver(), conn.Config.Driver, "the connection driver should be %s", unit.Driver())
}

func TestBuilderGetDB(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	db, err := builder.GetDB()
	assert.Equal(t, nil, err, "the return error should be nil")
	value := ""
	err = db.Get(&value, "SELECT 'hello' ")
	assert.Equal(t, nil, err, "the return error should be nil")
	if err == nil {
		assert.Equal(t, "hello", value, "the return value should be hello")
	}
}

func TestBuilderHasTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_builder")
	has, err := builder.HasTable("table_test_builder")
	assert.True(t, err == nil, "the return error should be nil")
	assert.False(t, has, "the return value should be false")
}

func TestBuilderCreateTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_builder")
	err := builder.CreateTable("table_test_builder", func(table Blueprint) {
		table.ID("id").Primary()
		table.UnsignedBigInteger("counter").Index()
		table.BigInteger("latest").Index()
		table.String("name", 20).Index()
		table.String("unionid", 128).Unique()
		table.AddUnique("name_latest", "name", "latest")
		table.AddIndex("name_counter", "name", "counter")
	})

	assert.True(t, builder.MustHasTable("table_test_builder"), "should return true")
	assert.Equal(t, nil, err, "the return error should be nil")

	// @todo: the return value should be refreshed
	// checkTable(t, table)
}

func TestBuilderGetTables(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	tables, err := builder.GetTables()
	assert.Equal(t, nil, err, "the return error should be nil")
	assert.True(t, utils.StringHave(tables, "table_test_builder"), "the return value should have table_test_builder")

	builder.DropTableIfExists("table_test_builder")
	tables, err = builder.GetTables()
	assert.Equal(t, nil, err, "the return error should be nil")
	assert.False(t, utils.StringHave(tables, "table_test_builder"), "the return value should have not table_test_builder")
}

func TestBuilderGetTable(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	table, err := builder.GetTable("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	assert.True(t, table != nil, "the return table should be BluePrint")
	if table == nil {
		return
	}
	checkTable(t, table)
}

func TestBuilderDropTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.DropTable("table_test_builder")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestBuilderDropTableIfExistsTableNotExists(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	err := builder.DropTableIfExists("table_not_exists")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestBuilderDropTableIfExistsTableExists(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	err := builder.DropTableIfExists("table_test_builder")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
	assert.Equal(t, nil, err, "the return error should be nil")
}

func TestBuilderRenameTable(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	err := builder.RenameTable("table_test_builder", "table_test_builder_re")
	assert.True(t, builder.MustHasTable("table_test_builder_re"), "should return true")
	assert.Equal(t, nil, err, "the return error should be nil")
	builder.DropTable("table_test_builder_re")
}

func TestBuilderAlterTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_builder")
	TestBuilderCreateTable(t)
	err := builder.AlterTable("table_test_builder", func(table Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid_unique")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
		table.AddIndex("nickname_index", "nickname")
		table.RenameIndex("latest_index", "re_latest_index")
	})
	assert.Equal(t, nil, err, "the return error should be nil")
	if err != nil {
		return
	}

	// @todo: the return value should be refreshed
	// checkTableAlterTable(t, table)

	// cheking the schema structure
	table, err := builder.GetTable("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	checkTableAlterTable(t, table)
	// builder.DropTable("table_test_builder")
}

func TestBuilderGetVersion(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	version, err := builder.GetVersion()
	assert.Equal(t, nil, err, "the return error should be nil")
	if err != nil {
		return
	}
	if unit.Is("mysql") {
		assert.Equal(t, "mysql", version.Driver, "the driver should be mysql")
		assert.Equal(t, 5, int(version.Major), "the major version should be 5")
		assert.Equal(t, 7, int(version.Minor), "the minor version should be 7")
	} else if unit.Is("mysql5.6") {
		assert.Equal(t, "mysql", version.Driver, "the driver should be mysql")
		assert.Equal(t, 5, int(version.Major), "the major version should be 5")
		assert.Equal(t, 6, int(version.Minor), "the minor version should be 6")
	} else if unit.Is("postgres") {
		assert.Equal(t, "postgres", version.Driver, "the driver should be postgres")
		assert.Equal(t, 9, int(version.Major), "the major version should be 9")
		assert.Equal(t, 6, int(version.Minor), "the minor version should be 6")
	} else if unit.Is("sqlite3") {
		assert.Equal(t, "sqlite3", version.Driver, "the driver should be sqlite3")
		assert.Equal(t, 3, int(version.Major), "the major version should be 3")
		assert.Equal(t, 34, int(version.Minor), "the minor version should be 34")
	}
	// fmt.Printf("The version is: %s %d.%d\n", version.Driver, version.Major, version.Minor)
}

func TestBuilderMustGetConnection(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	conn := builder.MustGetConnection()
	value := ""
	err := conn.DB.Get(&value, "SELECT 'hello' ")
	assert.Equal(t, nil, err, "the return error should be nil")
	if err == nil {
		assert.Equal(t, "hello", value, "the return value should be hello")
	}
	assert.Equal(t, unit.Driver(), conn.Config.Driver, "the connection driver should be %s", unit.Driver())
}

func TestBuilderMustGetDB(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	db := builder.MustGetDB()
	value := ""
	err := db.Get(&value, "SELECT 'hello' ")
	assert.Equal(t, nil, err, "the return error should be nil")
	if err == nil {
		assert.Equal(t, "hello", value, "the return value should be hello")
	}
}

func TestBuilderMustHasTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_builder")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
}

func TestBuilderMustGetTables(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	tables := builder.MustGetTables()
	assert.True(t, utils.StringHave(tables, "table_test_builder"), "the return value should have table_test_builder")

	builder.DropTableIfExists("table_test_builder")
	tables = builder.MustGetTables()
	assert.False(t, utils.StringHave(tables, "table_test_builder"), "the return value should have not table_test_builder")
}

func TestBuilderMustCreateTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.DropTableIfExists("table_test_builder")
	builder.MustCreateTable("table_test_builder", func(table Blueprint) {
		table.ID("id").Primary()
		table.UnsignedBigInteger("counter").Index()
		table.BigInteger("latest").Index()
		table.String("name", 20).Index()
		table.String("unionid", 128).Unique()
		table.AddUnique("name_latest", "name", "latest")
		table.AddIndex("name_counter", "name", "counter")
	})
	assert.True(t, builder.MustHasTable("table_test_builder"), "should return true")
}

func TestBuilderMustGetTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	table := builder.MustGetTable("table_test_builder")
	assert.True(t, table != nil, "the return table should be BluePrint")
	if table == nil {
		return
	}
	checkTable(t, table)
}

func TestBuilderMustDropTable(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.MustDropTable("table_test_builder")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
}

func TestBuilderMustDropTableIfExistsTableNotExists(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	builder.MustDropTableIfExists("table_not_exists")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
}

func TestBuilderMustDropTableIfExistsTableExists(t *testing.T) {
	defer unit.Catch()
	TestBuilderMustCreateTable(t)
	builder := getTestBuilder()
	builder.MustDropTableIfExists("table_test_builder")
	assert.False(t, builder.MustHasTable("table_test_builder"), "should return false")
}

func TestBuilderMustRenameTable(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	table := builder.MustRenameTable("table_test_builder", "table_test_builder_re")
	assert.True(t, builder.MustHasTable("table_test_builder_re"), "should return true")
	assert.Equal(t, "table_test_builder_re", table.GetName(), "the table name should be table_test_builder_re")
	builder.DropTable("table_test_builder_re")
}

func TestBuilderMustAlterTable(t *testing.T) {
	defer unit.Catch()
	TestBuilderCreateTable(t)
	builder := getTestBuilder()
	builder.MustAlterTable("table_test_builder", func(table Blueprint) {
		table.String("nickname", 50)
		table.String("unionid", 200)
		table.DropIndex("unionid_unique")
		table.DropColumn("name")
		table.RenameColumn("unionid", "uid").Unique()
		table.AddIndex("nickname_index", "nickname")
		table.RenameIndex("latest_index", "re_latest_index")
	})
	assert.True(t, builder.MustHasTable("table_test_builder"), "should return true")

	// cheking the schema structure
	table, err := builder.GetTable("table_test_builder")
	assert.Equal(t, nil, err, "the return error should be nil")
	checkTableAlterTable(t, table)
	builder.DropTable("table_test_builder")
}

func TestBuilderMustGetVersion(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	version := builder.MustGetVersion()
	if unit.Is("mysql") {
		assert.Equal(t, "mysql", version.Driver, "the driver should be mysql")
		assert.Equal(t, 5, int(version.Major), "the major version should be 5")
		assert.Equal(t, 7, int(version.Minor), "the minor version should be 7")
	} else if unit.Is("mysql5.6") {
		assert.Equal(t, "mysql", version.Driver, "the driver should be mysql")
		assert.Equal(t, 5, int(version.Major), "the major version should be 5")
		assert.Equal(t, 6, int(version.Minor), "the minor version should be 6")
	} else if unit.Is("postgres") {
		assert.Equal(t, "postgres", version.Driver, "the driver should be postgres")
		assert.Equal(t, 9, int(version.Major), "the major version should be 9")
		assert.Equal(t, 6, int(version.Minor), "the minor version should be 6")
	} else if unit.Is("sqlite3") {
		assert.Equal(t, "sqlite3", version.Driver, "the driver should be sqlite3")
		assert.Equal(t, 3, int(version.Major), "the major version should be 3")
		assert.Equal(t, 34, int(version.Minor), "the minor version should be 34")
	}
	// fmt.Printf("The version is: %s %d.%d\n", version.Driver, version.Major, version.Minor)
}

func TestBuilderDB(t *testing.T) {
	defer unit.Catch()
	builder := getTestBuilder()
	db := builder.DB()
	value := ""
	err := db.Get(&value, "SELECT 'hello' ")
	assert.Equal(t, nil, err, "the return error should be nil")
	if err == nil {
		assert.Equal(t, "hello", value, "the return value should be hello")
	}
}

// Utils..........

func checkTableAlterTable(t *testing.T, table Blueprint) {
	DefaultBigIntPrecision := 20
	if unit.Is("postgres") {
		DefaultBigIntPrecision = 64
	}

	columnType := "bigInteger"
	idColumnType := "bigInteger"
	if unit.Is("sqlite3") {
		idColumnType = "integer"
	}

	// checking the table schema structure
	assert.True(t, nil != table.GetColumn("id"), "the column id should be created")
	if table.GetColumn("id") != nil {
		assert.Equal(t, idColumnType, table.GetColumn("id").Type, "the id type should be %s", idColumnType)
		assert.Equal(t, "AutoIncrement", utils.StringVal(table.GetColumn("id").Extra), "the id extra should be AutoIncrement")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("id").Precision), "the id precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("id").IsUnsigned, "the id IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("counter"), "the column counter should be created")
	if table.GetColumn("counter") != nil {
		assert.Equal(t, columnType, table.GetColumn("counter").Type, "the counter type should be %s", columnType)
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("counter").Precision), "the counter precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("counter").IsUnsigned, "the counter IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("latest"), "the column latest should be created")
	if table.GetColumn("latest") != nil {
		assert.Equal(t, columnType, table.GetColumn("latest").Type, "the latest type should be %s", columnType)
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
		assert.Nil(t, nameLatest, "the index name_latest should has none")
	} else if unit.Not("sqlite3") {
		assert.Equal(t, 1, len(nameLatest.Columns), "the index name_latest should has one column")
		assert.Equal(t, "unique", nameLatest.Type, "the name_latest key type should be unique")
		if len(nameLatest.Columns) == 1 {
			assert.Equal(t, "latest", nameLatest.Columns[0].Name, "the second column of the index name_latest should be latest")
		}
	}

	nameCounter := table.GetIndex("name_counter")
	if unit.Is("postgres") {
		assert.Nil(t, nameCounter, "the index name_counter should has none")
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

	columnType := "bigInteger"
	idColumnType := "bigInteger"
	if unit.Is("sqlite3") {
		idColumnType = "integer"
	}

	// checking the table schema structure
	assert.True(t, nil != table.GetColumn("id"), "the column id should be created")
	if table.GetColumn("id") != nil {
		assert.Equal(t, idColumnType, table.GetColumn("id").Type, "the id type should be %s", idColumnType)
		assert.Equal(t, "AutoIncrement", utils.StringVal(table.GetColumn("id").Extra), "the id extra should be AutoIncrement")
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("id").Precision), "the id precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("id").IsUnsigned, "the id IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("counter"), "the column counter should be created")
	if table.GetColumn("counter") != nil {
		assert.Equal(t, columnType, table.GetColumn("counter").Type, "the counter type should be %s", columnType)
		assert.Equal(t, DefaultBigIntPrecision, utils.IntVal(table.GetColumn("counter").Precision), "the counter precision should be 20")
		if unit.Not("postgres") {
			assert.Equal(t, true, table.GetColumn("counter").IsUnsigned, "the counter IsUnsigned should be true")
		}
	}
	assert.True(t, nil != table.GetColumn("latest"), "the column latest should be created")
	if table.GetColumn("latest") != nil {
		assert.Equal(t, columnType, table.GetColumn("latest").Type, "the latest type should be %s", columnType)
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
