package postgres

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/jmoiron/sqlx"
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
// func (grammarSQL *Postgres) Confi2g(dsn string) {
// 	grammarSQL.DSN = dsn
// 	uinfo, err := url.Parse(grammarSQL.DSN)
// 	if err != nil {
// 		panic(err)
// 	}
// 	grammarSQL.DB = filepath.Base(uinfo.Path)
// 	schema := uinfo.Query().Get("search_path")
// 	if schema == "" {
// 		schema = "public"
// 	}
// 	grammarSQL.Schema = schema
// }

// setup the method will be executed when db server was connected
func (grammarSQL *Postgres) setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {
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
	grammarSQL.DatabaseName = filepath.Base(uinfo.Path)
	schema := uinfo.Query().Get("search_path")
	if schema == "" {
		schema = "public"
	}
	grammarSQL.SchemaName = schema
	return nil
}

// NewWith Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL Postgres) NewWith(db *sqlx.DB, config *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
	err := grammarSQL.setup(db, config, option)
	if err != nil {
		return nil, err
	}
	grammarSQL.Quoter.Bind(db, option.Prefix)
	return grammarSQL, nil
}

// NewWithRead Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL Postgres) NewWithRead(write *sqlx.DB, writeConfig *dbal.Config, read *sqlx.DB, readConfig *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
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
func New() dbal.Grammar {
	pg := Postgres{
		SQL: sql.NewSQL(&Quoter{}),
	}
	pg.Driver = "postgres"
	pg.IndexTypes = map[string]string{
		"unique": "UNIQUE INDEX",
		"index":  "INDEX",
	}

	// overwrite types
	types := pg.SQL.Types
	types["tinyInteger"] = "SMALLINT"
	types["bigInteger"] = "BIGINT"
	types["string"] = "CHARACTER VARYING"
	types["integer"] = "INTEGER"
	types["decimal"] = "NUMERIC"
	types["float"] = "REAL"
	types["double"] = "DOUBLE PRECISION"
	types["char"] = "CHARACTER"
	types["mediumText"] = "TEXT"
	types["longText"] = "TEXT"
	types["dateTime"] = "TIMESTAMP(%d) WITHOUT TIME ZONE"
	types["dateTimeTz"] = "TIMESTAMP(%d) WITH TIME ZONE"
	types["time"] = "TIME(%d) WITHOUT TIME ZONE"
	types["timeTz"] = "TIME(%d) WITH TIME ZONE"
	types["timestamp"] = "TIMESTAMP(%d) WITHOUT TIME ZONE"
	types["timestampTz"] = "TIMESTAMP(%d) WITH TIME ZONE"
	types["binary"] = "BYTEA"
	types["macAddress"] = "MACADDR"
	pg.Types = types

	// set fliptypes
	flipTypes, ok := utils.MapFilp(pg.Types)
	if ok {
		pg.FlipTypes = flipTypes.(map[string]string)
		pg.FlipTypes["TEXT"] = "text"
		pg.FlipTypes["TIMESTAMP WITHOUT TIME ZONE"] = "timestamp"
		pg.FlipTypes["TIMESTAMP WITH TIME ZONE"] = "timestampTz"
		pg.FlipTypes["TIME WITHOUT TIME ZONE"] = "time"
		pg.FlipTypes["TIME WITH TIME ZONE"] = "timeTz"
		pg.FlipTypes["SMALLINT"] = "smallInteger"
	}

	return pg
}

// GetOperators get the operators
func (grammarSQL Postgres) GetOperators() []string {
	return []string{
		"=", "<", ">", "<=", ">=", "<>", "!=",
		"like", "not like", "between", "ilike", "not ilike",
		"~", "&", "|", "#", "<<", ">>", "<<=", ">>=",
		"&&", "@>", "<@", "?", "?|", "?&", "||", "-", "@?", "@@", "#-",
		"is distinct from", "is not distinct from",
	}
}
