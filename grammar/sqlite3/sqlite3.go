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
func (grammarSQL *SQLite3) setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {

	if db == nil {
		return fmt.Errorf("db is nil")
	}

	if config == nil {
		return fmt.Errorf("config is nil")
	}

	grammarSQL.DB = db
	grammarSQL.Config = config
	grammarSQL.Option = option
	uinfo, err := url.Parse(grammarSQL.Config.DSN)
	if err != nil {
		return err
	}
	filename := filepath.Base(uinfo.Path)
	grammarSQL.DatabaseName = strings.TrimSuffix(filename, filepath.Ext(filename))
	grammarSQL.SchemaName = grammarSQL.DatabaseName
	return nil
}

// NewWith Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL SQLite3) NewWith(db *sqlx.DB, config *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
	err := grammarSQL.setup(db, config, option)
	if err != nil {
		return nil, err
	}
	grammarSQL.Quoter.Bind(db, option.Prefix)
	return grammarSQL, nil
}

// NewWithRead Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL SQLite3) NewWithRead(write *sqlx.DB, writeConfig *dbal.Config, read *sqlx.DB, readConfig *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
	err := grammarSQL.setup(write, writeConfig, option)
	if err != nil {
		return nil, err
	}

	grammarSQL.Read = read
	grammarSQL.ReadConfig = readConfig
	grammarSQL.Quoter.Bind(write, option.Prefix, read)
	return grammarSQL, nil
}

// New Create a new mysql grammar inteface
func New(opts ...sql.Option) dbal.Grammar {
	sqlite := SQLite3{
		SQL: sql.NewSQL(&Quoter{}, opts...),
	}
	if sqlite.Driver == "" || sqlite.Driver == "sql" {
		sqlite.Driver = "sqlite3"
	}
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

	return sqlite
}

// GetOperators get the operators
func (grammarSQL SQLite3) GetOperators() []string {
	return []string{
		"=", "<", ">", "<=", ">=", "<>", "!=",
		"like", "not like", "ilike",
		"&", "|", "<<", ">>",
	}
}
