package sqlite3

import (
	"net/url"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// SQLite3 the sqlite3 Grammar
type SQLite3 struct {
	sql.SQL
}

func init() {
	dbal.Register("sqlite3", New())
}

// Config set the configure using DSN
func (grammarSQL *SQLite3) Config(dsn string) {
	grammarSQL.DSN = dsn
	uinfo, err := url.Parse(grammarSQL.DSN)
	if err != nil {
		grammarSQL.DB = "memory"
	} else {
		filename := filepath.Base(uinfo.Path)
		grammarSQL.DB = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	grammarSQL.Schema = grammarSQL.DB
}

// New Create a new mysql grammar inteface
func New() dbal.Grammar {
	sqlite := SQLite3{
		SQL: sql.NewSQL(sql.Quoter{}),
	}
	sqlite.Driver = "sqlite3"
	sqlite.IndexTypes = map[string]string{
		"unique": "UNIQUE INDEX",
		"index":  "INDEX",
	}

	// overwrite types
	sqlite.Types["bigInteger"] = "INTEGER"
	sqlite.Types["smallInteger"] = "INTEGER"
	sqlite.Types["integer"] = "INTEGER"
	sqlite.Types["char"] = "CHARACTER"

	// set fliptypes
	flipTypes, ok := utils.MapFilp(sqlite.Types)
	if ok {
		sqlite.FlipTypes = flipTypes.(map[string]string)
		sqlite.FlipTypes["INTEGER"] = "integer"
		sqlite.FlipTypes["DATETIME"] = "dateTime"
		sqlite.FlipTypes["TIME"] = "time"
	}

	return &sqlite
}
