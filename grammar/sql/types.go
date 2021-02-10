package sql

import (
	"github.com/yaoapp/xun/grammar"
)

// SQL the SQL Grammar
type SQL struct {
	Driver     string
	Mode       string
	Types      map[string]string
	IndexTypes map[string]string
	Quoter     grammar.Quoter
	Builder    grammar.SQLBuilder
}

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	sql := NewSQL()
	return &sql
}

// NewSQL create a new SQL instance
func NewSQL() SQL {
	return SQL{
		Driver:  "sql",
		Mode:    "production",
		Quoter:  Quoter{},
		Builder: Builder{},
		IndexTypes: map[string]string{
			"primary": "PRIMARY KEY",
			"unique":  "UNIQUE KEY",
			"index":   "KEY",
		},
		Types: map[string]string{
			"bigIncrements":         "VARCHAR",
			"bigInteger":            "VARCHAR",
			"binary":                "VARCHAR",
			"boolean":               "VARCHAR",
			"char":                  "VARCHAR",
			"dateTimeTz":            "VARCHAR",
			"dateTime":              "VARCHAR",
			"date":                  "VARCHAR",
			"decimal":               "VARCHAR",
			"double":                "VARCHAR",
			"enum":                  "VARCHAR",
			"float":                 "VARCHAR",
			"foreignId":             "VARCHAR",
			"geometryCollection":    "VARCHAR",
			"geometry":              "VARCHAR",
			"id":                    "VARCHAR",
			"increments":            "VARCHAR",
			"integer":               "VARCHAR",
			"ipAddress":             "VARCHAR",
			"json":                  "VARCHAR",
			"jsonb":                 "VARCHAR",
			"lineString":            "VARCHAR",
			"longText":              "VARCHAR",
			"macAddress":            "VARCHAR",
			"mediumIncrements":      "VARCHAR",
			"mediumInteger":         "VARCHAR",
			"mediumText":            "VARCHAR",
			"morphs":                "VARCHAR",
			"multiLineString":       "VARCHAR",
			"multiPoint":            "VARCHAR",
			"multiPolygon":          "VARCHAR",
			"nullableMorphs":        "VARCHAR",
			"nullableTimestamps":    "VARCHAR",
			"nullableUuidMorphs":    "VARCHAR",
			"point":                 "VARCHAR",
			"polygon":               "VARCHAR",
			"rememberToken":         "VARCHAR",
			"set":                   "VARCHAR",
			"smallIncrements":       "VARCHAR",
			"smallInteger":          "VARCHAR",
			"softDeletesTz":         "VARCHAR",
			"softDeletes":           "VARCHAR",
			"string":                "VARCHAR",
			"text":                  "VARCHAR",
			"timeTz":                "VARCHAR",
			"time":                  "VARCHAR",
			"timestampTz":           "VARCHAR",
			"timestamp":             "VARCHAR",
			"timestampsTz":          "VARCHAR",
			"timestamps":            "VARCHAR",
			"tinyIncrements":        "VARCHAR",
			"tinyInteger":           "VARCHAR",
			"unsignedBigInteger":    "VARCHAR",
			"unsignedDecimal":       "VARCHAR",
			"unsignedInteger":       "VARCHAR",
			"unsignedMediumInteger": "VARCHAR",
			"unsignedSmallInteger":  "VARCHAR",
			"unsignedTinyInteger":   "VARCHAR",
			"uuidMorphs":            "VARCHAR",
			"uuid":                  "VARCHAR",
			"year":                  "VARCHAR",
		},
	}
}
