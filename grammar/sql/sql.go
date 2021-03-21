package sql

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/utils"
)

// SQL the SQL Grammar
type SQL struct {
	Driver       string
	Mode         string
	Types        map[string]string
	FlipTypes    map[string]string
	IndexTypes   map[string]string
	DatabaseName string
	SchemaName   string
	DB           *sqlx.DB
	Config       *dbal.Config
	Option       *dbal.Option
	dbal.Grammar
	dbal.Quoter
}

// NewSQL create a new SQL instance
func NewSQL(quoter dbal.Quoter) SQL {
	sql := &SQL{
		Driver: "sql",
		Mode:   "production",
		Quoter: quoter,
		IndexTypes: map[string]string{
			"unique": "UNIQUE KEY",
			"index":  "KEY",
		},
		FlipTypes: map[string]string{},
		Types: map[string]string{
			"tinyInteger":  "TINYINT",
			"smallInteger": "SMALLINT",
			"integer":      "INT",
			"bigInteger":   "BIGINT",
			"boolean":      "BOOLEAN",
			"decimal":      "DECIMAL",
			"float":        "FLOAT",
			"double":       "DOUBLE",
			"string":       "VARCHAR",
			"char":         "CHAR",
			"text":         "TEXT",
			"mediumText":   "MEDIUMTEXT",
			"longText":     "LONGTEXT",
			"binary":       "VARBINARY",
			"date":         "DATE",
			"dateTime":     "DATETIME",
			"dateTimeTz":   "DATETIME",
			"time":         "TIME",
			"timeTz":       "TIME",
			"timestamp":    "TIMESTAMP",
			"timestampTz":  "TIMESTAMP",
			"enum":         "ENUM",
			"json":         "JSON",
			"jsonb":        "JSONB",
			"uuid":         "UUID",
			"ipAddress":    "IPADDRESS",
			"macAddress":   "MACADDRESS",
			"year":         "YEAR",
			// "mediumInteger": "mediumInteger",
		},
	}
	return *sql
}

// New Create a new mysql grammar inteface
func New(dsn string) dbal.Grammar {
	sql := NewSQL(Quoter{})
	flipTypes, ok := utils.MapFilp(sql.Types)
	if ok {
		sql.FlipTypes = flipTypes.(map[string]string)
	}
	return &sql
}

// Setup the method will be executed when db server was connected
func (grammarSQL *SQL) Setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {

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
	grammarSQL.SchemaName = grammarSQL.DatabaseName
	return nil
}

// OnConnected the event will be triggered when db server was connected
func (grammarSQL *SQL) OnConnected() error {
	return nil
}
