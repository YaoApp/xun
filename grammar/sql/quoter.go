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
func (quoter *Quoter) ID(value string) string {
	value = strings.ReplaceAll(value, "`", "")
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	return "`" + value + "`"
}

// VAL quoting query value ( 'value' )
func (quoter *Quoter) VAL(value interface{}) string {
	input := ""
	switch value.(type) {
	case *string:
		input = fmt.Sprintf("%s", utils.StringVal(value.(*string)))
		break
	case string:
		input = fmt.Sprintf("%s", value)
		break
	case int, int16, int32, int64, float64, float32:
		input = fmt.Sprintf("%d", value)
		break
	default:
		input = fmt.Sprintf("%v", value)
	}

	input = strings.ReplaceAll(input, "'", "\\'")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return "'" + input + "'"
}

// Wrap a value in keyword identifiers.
func (quoter *Quoter) Wrap(value interface{}) string {
	switch value.(type) {
	case dbal.Expression:
		return value.(dbal.Expression).GetValue()
	case dbal.Name:
		return quoter.ID(value.(dbal.Name).Fullname())
	case string:
		return quoter.ID(dbal.NewName(value.(string)).Fullname())
	default:
		return fmt.Sprintf("%v", value)
	}
}

// IsExpression Determine if the given value is a raw expression.
func (quoter *Quoter) IsExpression(value interface{}) bool {
	switch value.(type) {
	case dbal.Expression:
		return true
	default:
		return false
	}
}

// Parameter Get the appropriate query parameter place-holder for a value.
func (quoter *Quoter) Parameter(value interface{}, num int) string {
	if quoter.IsExpression(value) {
		return value.(dbal.Expression).GetValue()
	}
	return "?"
}

// Parameterize Create query parameter place-holders for an array.
func (quoter *Quoter) Parameterize(values []interface{}, offset int) string {
	params := []string{}
	for idx, value := range values {
		params = append(params, quoter.Parameter(value, offset+idx+1))
	}
	return strings.Join(params, ",")
}

// Columnize Convert an array of column names into a delimited string.
func (quoter *Quoter) Columnize(columns []interface{}) string {
	wrapColumns := []string{}
	for _, col := range columns {
		switch col.(type) {
		case dbal.Name:
			wrapColumns = append(wrapColumns, quoter.ID(col.(dbal.Name).Name))
		case dbal.Expression:
			wrapColumns = append(wrapColumns, col.(dbal.Expression).GetValue())
		}
	}
	return strings.Join(wrapColumns, ", ")
}
