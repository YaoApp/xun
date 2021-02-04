package query

import (
	"testing"
)

func TestTableMySQL(t *testing.T) {
	builder := TableMySQL()
	builder.Where()
	builder.Join()
}

func TestTableSQLite(t *testing.T) {
	builder := TableSQLite()
	builder.Where()
	builder.Join()
}

func TestTableSQLServer(t *testing.T) {
	builder := TableSQLServer()
	builder.Where()
	builder.Join()
}

func TestTableOracle(t *testing.T) {
	builder := TableOracle()
	builder.Where()
	builder.Join()
}
