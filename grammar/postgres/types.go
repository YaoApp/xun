package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// Postgres the Postgresql Grammar
type Postgres struct {
	sql.SQL
}

// Quoter the database quoting query text SQL type
type Quoter struct{}

// ID quoting query Identifier (`id`)
func (quoter Quoter) ID(name string, db *sqlx.DB) string {
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	return "\"" + name + "\""
}

// VAL quoting query value ( 'value' )
func (quoter Quoter) VAL(v interface{}, db *sqlx.DB) string {
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

// New Create a new mysql grammar inteface
func New(dsn string) grammar.Grammar {
	pg := Postgres{
		SQL: sql.NewSQL(dsn, Quoter{}),
	}
	pg.Driver = "postgres"
	pg.IndexTypes = map[string]string{
		"primary": "PRIMARY KEY",
		"unique":  "UNIQUE INDEX",
		"index":   "INDEX",
	}

	// update schema name
	pg.SchemaName()

	// overwrite types
	types := pg.SQL.Types
	types["bigInteger"] = "BIGINT"
	types["string"] = "CHARACTER VARYING"
	pg.Types = types

	// set fliptypes
	flipTypes, ok := utils.MapFilp(pg.Types)
	if ok {
		pg.FlipTypes = flipTypes.(map[string]string)
		// pg.FlipTypes["CHARACTER VARYING"] = "string"
	}
	return &pg
}
