package mysql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// Mysql the mysql Grammar
type Mysql struct {
	sql.SQL
}

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

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	my := Mysql{
		SQL: sql.NewSQL(),
	}
	my.Driver = "mysql"
	my.Quoter = Quoter{}

	// set fliptypes
	flipTypes, ok := utils.MapFilp(my.Types)
	if ok {
		my.FlipTypes = flipTypes.(map[string]string)
	}
	return &my
}
