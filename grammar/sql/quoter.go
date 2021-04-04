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
	Prefix     string
}

//Bind make a new Quoter inteface
func (quoter *Quoter) Bind(db *sqlx.DB, prefix string, dbRead ...*sqlx.DB) dbal.Quoter {
	quoter.DBPrimary = db
	quoter.DB = db
	quoter.Prefix = prefix
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
		col := value.(dbal.Name)
		if col.As() != "" {
			return fmt.Sprintf("%s as %s", quoter.ID(col.Name), col.As())
		}
		return quoter.ID(value.(dbal.Name).Name)
	case string:
		str := value.(string)
		if str == "*" {
			return str
		}
		if strings.Contains(str, ".") {
			arrs := strings.Split(str, ".")
			tab := arrs[0]
			col := dbal.NewName(arrs[1])
			if col.As() != "" {
				return fmt.Sprintf("%s.%s as %s", quoter.ID(tab), quoter.ID(col.Name), quoter.ID(col.As()))
			}
			name := col.Name
			if name != "*" {
				name = quoter.ID(col.Name)
			}
			return fmt.Sprintf("%s.%s", quoter.ID(tab), name)
		}
		return quoter.ID(dbal.NewName(value.(string)).Name)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// WrapTable Wrap a table in keyword identifiers.
func (quoter *Quoter) WrapTable(value interface{}) string {
	switch value.(type) {
	case dbal.Expression:
		return value.(dbal.Expression).GetValue()
	case dbal.Name:
		col := value.(dbal.Name)
		if col.As() != "" {
			return fmt.Sprintf("%s as %s", col.Fullname(), col.As())
		}
		return quoter.ID(value.(dbal.Name).Fullname())
	case string:
		return quoter.ID(dbal.NewName(value.(string)).Fullname())
	default:
		return fmt.Sprintf("%v", value)
	}
}

// WrapUnion a union subquery in parentheses.
func (quoter *Quoter) WrapUnion(sql string) string {
	return fmt.Sprintf("(%s)", sql)
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
		wrapColumns = append(wrapColumns, quoter.Wrap(col))
	}
	return strings.Join(wrapColumns, ", ")
}
