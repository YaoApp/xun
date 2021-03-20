package sqlite3

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
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

// Config2 set the configure using DSN
// func (grammarSQL *SQLite3) Config2(dsn string) {
// 	grammarSQL.DSN = dsn
// 	uinfo, err := url.Parse(grammarSQL.DSN)
// 	if err != nil {
// 		grammarSQL.DB = "memory"
// 	} else {
// 		filename := filepath.Base(uinfo.Path)
// 		grammarSQL.DB = strings.TrimSuffix(filename, filepath.Ext(filename))
// 	}
// 	grammarSQL.Schema = grammarSQL.DB
// }

// Setup the method will be executed when db server was connected
func (grammarSQL *SQLite3) Setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {
	grammarSQL.DB = db
	grammarSQL.Config = config
	grammarSQL.Option = option
	if grammarSQL.Config == nil {
		return fmt.Errorf("config is nil")
	}
	uinfo, err := url.Parse(grammarSQL.Config.DSN)
	if err != nil {
		return err
	}
	filename := filepath.Base(uinfo.Path)
	grammarSQL.DatabaseName = strings.TrimSuffix(filename, filepath.Ext(filename))
	grammarSQL.SchemaName = grammarSQL.DatabaseName
	return nil
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
	sqlite.Types["tinyInteger"] = "TINYINT"
	sqlite.Types["bigInteger"] = "BIGINT"
	sqlite.Types["smallInteger"] = "SMALLINT"
	sqlite.Types["integer"] = "INTEGER"
	sqlite.Types["char"] = "CHARACTER"
	sqlite.Types["binary"] = "BLOB"

	// set fliptypes
	flipTypes, ok := utils.MapFilp(sqlite.Types)
	if ok {
		sqlite.FlipTypes = flipTypes.(map[string]string)
		sqlite.FlipTypes["DATETIME"] = "dateTime"
		sqlite.FlipTypes["TIME"] = "time"
		sqlite.FlipTypes["TIMESTAMP"] = "timestamp"
		sqlite.FlipTypes["UNSIGNED BIG INT"] = "bigInteger"
	}

	return &sqlite
}
