package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/dbal/schema"
	"github.com/yaoapp/xun/unit"
)

var testQueryBuilder Query
var testQueryBuilderInstance *Builder
var testSchemaBuilder schema.Schema

func getTestBuilder() Query {
	defer unit.Catch()
	unit.SetLogger()
	if testQueryBuilder != nil {
		return testQueryBuilder
	}
	testQueryBuilder = New(unit.Driver(), unit.DSN())
	return testQueryBuilder
}

func getTestBuilderInstance() *Builder {
	defer unit.Catch()
	unit.SetLogger()
	if testQueryBuilderInstance != nil {
		return testQueryBuilderInstance
	}
	testQueryBuilderInstance = newBuilder(unit.Driver(), unit.DSN())
	return testQueryBuilderInstance
}

func getTestSchemaBuilder() schema.Schema {
	defer unit.Catch()
	unit.SetLogger()
	if testSchemaBuilder != nil {
		return testSchemaBuilder
	}
	driver := unit.Driver()
	dsn := unit.DSN()
	testSchemaBuilder = schema.New(driver, dsn)
	return testSchemaBuilder
}

func TestBuilderNew(t *testing.T) {
	defer unit.Catch()
	unit.SetLogger()
	qb := getTestBuilder()
	db := qb.DB()
	assert.True(t, db.Ping() == nil, "The primary connection should be made")

	dbReadonly := qb.DB(true)
	assert.True(t, dbReadonly.Ping() == nil, "The read-only connection should be made")
}

func TestBuilderDriver(t *testing.T) {
	defer unit.Catch()
	unit.SetLogger()
	qb := getTestBuilder()
	driver, err := qb.Driver()
	assert.Nil(t, err, "The error should be nil")
	assert.Equal(t, unit.Driver(), driver, "The driver should be mysql")
}
