package sqlite3

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// SQLite3 the sqlite3 Grammar
type SQLite3 struct {
	sql.SQL
}

// New Create a new mysql grammar inteface
func New(dsn string) grammar.Grammar {
	sqlite := SQLite3{
		SQL: sql.NewSQL(dsn, sql.Quoter{}),
	}
	sqlite.Driver = "sqlite3"
	sqlite.IndexTypes = map[string]string{
		"unique": "UNIQUE INDEX",
		"index":  "INDEX",
	}

	// overwrite types
	types := sqlite.SQL.Types
	types["bigInteger"] = "INTEGER"
	sqlite.Types = types

	// set fliptypes
	flipTypes, ok := utils.MapFilp(sqlite.Types)
	if ok {
		sqlite.FlipTypes = flipTypes.(map[string]string)
	}
	return &sqlite
}
