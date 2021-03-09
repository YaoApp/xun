package postgres

import (
	"net/url"
	"path/filepath"

	_ "github.com/lib/pq" // Load postgres driver
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// Postgres the Postgresql Grammar
type Postgres struct {
	sql.SQL
}

func init() {
	dbal.Register("postgres", New())
}

// Config set the configure using DSN
func (grammarSQL *Postgres) Config(dsn string) {
	grammarSQL.DSN = dsn
	uinfo, err := url.Parse(grammarSQL.DSN)
	if err != nil {
		panic(err)
	}
	grammarSQL.DB = filepath.Base(uinfo.Path)
	schema := uinfo.Query().Get("search_path")
	if schema == "" {
		schema = "public"
	}
	grammarSQL.Schema = schema
}

// New Create a new mysql grammar inteface
func New() dbal.Grammar {
	pg := Postgres{
		SQL: sql.NewSQL(Quoter{}),
	}
	pg.Driver = "postgres"
	pg.IndexTypes = map[string]string{
		"unique": "UNIQUE INDEX",
		"index":  "INDEX",
	}

	// overwrite types
	types := pg.SQL.Types
	types["bigInteger"] = "BIGINT"
	types["string"] = "CHARACTER VARYING"
	types["integer"] = "INTEGER"
	types["decimal"] = "NUMERIC"
	types["float"] = "REAL"
	pg.Types = types

	// set fliptypes
	flipTypes, ok := utils.MapFilp(pg.Types)
	if ok {
		pg.FlipTypes = flipTypes.(map[string]string)
		// pg.FlipTypes["CHARACTER VARYING"] = "string"
	}
	return &pg
}
