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
			return fmt.Sprintf("%s as %s", quoter.ID(col.Name), col.As())
		}
		return quoter.ID(value.(dbal.Name).Name)
	case dbal.Select:
		col := value.(dbal.Select)
		if col.Alias != "" {
			return fmt.Sprintf("%s as %s", col.SQL, quoter.ID(col.Alias))
		}
		return fmt.Sprintf("%s ", col.SQL)
	case string:
		return quoter.WrapAliasedValue(value.(string))
	default:
		return fmt.Sprintf("%v", value)
	}
}

// WrapAliasedValue Wrap a value that has an alias.
func (quoter *Quoter) WrapAliasedValue(value string) string {
	if value == "*" {
		return "*"
	}
	if strings.Contains(value, ".") {
		arrs := strings.Split(value, ".")
		table := arrs[0]
		name := quoter.WrapAliasedValue(arrs[1])
		return fmt.Sprintf("%s.%s", quoter.ID(table), name)
	}

	name := dbal.NewName(value)
	if name.As() != "" {
		return fmt.Sprintf("%s as %s", quoter.ID(name.Fullname()), quoter.ID(name.As()))
	}
	return fmt.Sprintf("%s", quoter.ID(name.Fullname()))
}

// WrapTable Wrap a table in keyword identifiers.
func (quoter *Quoter) WrapTable(value interface{}) string {
	switch value.(type) {
	case dbal.Expression:
		return value.(dbal.Expression).GetValue()
	case dbal.Name:
		col := value.(dbal.Name)
		if col.As() != "" {
			return fmt.Sprintf("%s as %s", quoter.ID(col.Fullname()), quoter.ID(col.As()))
		}
		return quoter.ID(value.(dbal.Name).Fullname())
	case dbal.From:
		return quoter.WrapTable(value.(dbal.From).Name)
	case string:
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
		params = append(params, quoter.Parameter(value, idx+1+offset))
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
