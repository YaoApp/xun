package sqlite3

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
)

// SQLite3 the sqlite3 Grammar
type SQLite3 struct {
	sql.SQL
}

// Builder sqlite SQL builder
type Builder struct {
	sql.Builder
}

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	sqlite := SQLite3{
		SQL: sql.NewSQL(),
	}
	sqlite.Driver = "sqlite3"
	sqlite.Builder = Builder{}
	sqlite.IndexTypes = map[string]string{
		"unique": "UNIQUE INDEX",
		"index":  "INDEX",
	}
	types := sqlite.SQL.Types
	types["bigInteger"] = "INTEGER"
	sqlite.Types = types
	return &sqlite
}
