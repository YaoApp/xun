package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Quoter the database quoting query text SQL type
type Quoter struct{}

// ID quoting query Identifier (`id`)
func (quoter Quoter) ID(name string, db *sqlx.DB) string {
	name = strings.ReplaceAll(name, "`", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	return "`" + name + "`"
}

// VAL quoting query value ( 'value' )
func (quoter Quoter) VAL(v interface{}, db *sqlx.DB) string {
	input := fmt.Sprintf("%v", v)
	input = strings.ReplaceAll(input, "'", "\\'")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return "'" + input + "'"
}
