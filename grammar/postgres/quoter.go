package postgres

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// Quoter the database quoting query text SQL type
type Quoter struct {
	sql.Quoter
}

// ID quoting query Identifier (`id`)
func (quoter Quoter) ID(name string) string {
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	return "\"" + name + "\""
}

// VAL quoting query value ( 'value' )
func (quoter Quoter) VAL(v interface{}) string {
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

// Wrap a value in keyword identifiers.
func (quoter *Quoter) Wrap(value interface{}) string {
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
		str := value.(string)
		if strings.Contains(str, ".") {
			arrs := strings.Split(str, ".")
			tab := arrs[0]
			col := dbal.NewName(arrs[1])
			if col.As() != "" {
				return fmt.Sprintf("%s.%s as %s", quoter.ID(tab), quoter.ID(col.Fullname()), quoter.ID(col.As()))
			}
			return fmt.Sprintf("%s.%s", quoter.ID(tab), quoter.ID(col.Fullname()))
		}
		return quoter.ID(dbal.NewName(value.(string)).Fullname())
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Parameter Get the appropriate query parameter place-holder for a value.
func (quoter *Quoter) Parameter(value interface{}, num int) string {
	if quoter.IsExpression(value) {
		return value.(dbal.Expression).GetValue()
	}
	return fmt.Sprintf("$%d", num)
}

// Parameterize Create query parameter place-holders for an array.
func (quoter *Quoter) Parameterize(values []interface{}, offset int) string {
	params := []string{}
	for idx, value := range values {
		quoter.Parameter(value, idx+1+offset)
	}
	return strings.Join(params, ",")
}
