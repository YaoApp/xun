package sql

import (
	"database/sql"
	"fmt"

	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/dbal"
)

// Upsert Upsert new records or update the existing ones.
func (grammarSQL SQL) Upsert(query *dbal.Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (sql.Result, error) {
	return nil, fmt.Errorf("This database engine does not support upserts")
}

// CompileUpsert Upsert new records or update the existing ones.
func (grammarSQL SQL) CompileUpsert(query *dbal.Query, values []xun.R, uniqueBy []interface{}, updateValues interface{}) (string, []interface{}) {
	panic(fmt.Errorf("This database engine does not support upserts"))
}
