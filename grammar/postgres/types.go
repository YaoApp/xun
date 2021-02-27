package postgres

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// Postgres the Postgresql Grammar
type Postgres struct {
	sql.SQL
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
