package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// Quoter the database quoting query text SQL type
type Quoter struct {
	DB         *sqlx.DB
	DBPrimary  *sqlx.DB
	DBReadOnly *sqlx.DB
}

//Bind make a new Quoter inteface
func (quoter *Quoter) Bind(db *sqlx.DB, dbRead ...*sqlx.DB) dbal.Quoter {
	quoter.DBPrimary = db
	quoter.DB = db
	if len(dbRead) > 0 && dbRead[0] != nil {
		quoter.DBReadOnly = dbRead[0]
	}
	return quoter
}

// Read pick the readonly connection
func (quoter *Quoter) Read() dbal.Quoter {
	quoter.DB = quoter.DBReadOnly
	return quoter
}

// Write pick the primary connection
func (quoter *Quoter) Write() dbal.Quoter {
	quoter.DB = quoter.DBPrimary
	return quoter
}

// ID quoting query Identifier (`id`)
func (quoter *Quoter) ID(name string, db *sqlx.DB) string {
	name = strings.ReplaceAll(name, "`", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	return "`" + name + "`"
}

// VAL quoting query value ( 'value' )
func (quoter *Quoter) VAL(v interface{}, db *sqlx.DB) string {
	input := ""
	switch v.(type) {
	case *string:
		input = fmt.Sprintf("%s", utils.StringVal(v.(*string)))
		break
	case string:
		input = fmt.Sprintf("%s", v)
		break
	case int, int16, int32, int64, float64, float32:
		input = fmt.Sprintf("%d", v)
		break
	default:
		input = fmt.Sprintf("%v", v)
	}

	input = strings.ReplaceAll(input, "'", "\\'")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return "'" + input + "'"
}
