package sql

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/utils"
)

// SQL the SQL Grammar
type SQL struct {
	Driver     string
	Mode       string
	Types      map[string]string
	FlipTypes  map[string]string
	IndexTypes map[string]string
	Quoter     grammar.Quoter
	DSN        string
	DB         string
	Schema     string
	grammar.Grammar
}

// New Create a new mysql grammar inteface
func New(dsn string) grammar.Grammar {
	sql := NewSQL(dsn, Quoter{})
	flipTypes, ok := utils.MapFilp(sql.Types)
	if ok {
		sql.FlipTypes = flipTypes.(map[string]string)
	}
	return &sql
}

// NewSQL create a new SQL instance
func NewSQL(dsn string, quoter grammar.Quoter) SQL {
	sql := &SQL{
		Driver: "sql",
		DSN:    dsn,
		Mode:   "production",
		Quoter: quoter,
		IndexTypes: map[string]string{
			"primary": "PRIMARY KEY",
			"unique":  "UNIQUE KEY",
			"index":   "KEY",
		},
		FlipTypes: map[string]string{},
		Types: map[string]string{
			"bigInteger":    "BIGINT",
			"string":        "VARCHAR",
			"binary":        "binary",
			"boolean":       "boolean",
			"char":          "char",
			"dateTimeTz":    "dateTimeTz",
			"dateTime":      "dateTime",
			"date":          "date",
			"decimal":       "decimal",
			"double":        "double",
			"enum":          "enum",
			"float":         "float",
			"integer":       "integer",
			"json":          "JSON",
			"jsonb":         "JSONB",
			"longText":      "LONGTEXT",
			"mediumInteger": "mediumInteger",
			"mediumText":    "mediumText",
			"smallInteger":  "smallInteger",
			"text":          "text",
			"timestamp":     "timestamp",
			"timestampsTz":  "timestampsTz",
			"tinyInteger":   "tinyInteger",
			"uuid":          "UUID",
			"year":          "YEAR",
		},
	}
	sql.DBName()
	sql.SchemaName()
	return *sql
}
