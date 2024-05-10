package mysql

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //Load mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/dbal"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// MySQL the MySQL Grammar
type MySQL struct {
	sql.SQL
}

func init() {
	dbal.Register("mysql", New())
}

// setup the method will be executed when db server was connected
func (grammarSQL *MySQL) setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	if config == nil {
		return fmt.Errorf("config is nil")
	}

	grammarSQL.DB = db
	grammarSQL.Config = config
	grammarSQL.Option = option
	cfg, err := mysql.ParseDSN(grammarSQL.Config.DSN)
	if err != nil {
		return err
	}
	grammarSQL.DatabaseName = cfg.DBName
	grammarSQL.SchemaName = grammarSQL.DatabaseName
	return nil
}

// NewWith Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL MySQL) NewWith(db *sqlx.DB, config *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
	err := grammarSQL.setup(db, config, option)
	if err != nil {
		return nil, err
	}
	grammarSQL.Quoter.Bind(db, option.Prefix)
	return grammarSQL, nil
}

// NewWithRead Create a new grammar interface, using the given *sqlx.DB, *dbal.Config and *dbal.Option.
func (grammarSQL MySQL) NewWithRead(write *sqlx.DB, writeConfig *dbal.Config, read *sqlx.DB, readConfig *dbal.Config, option *dbal.Option) (dbal.Grammar, error) {
	err := grammarSQL.setup(write, writeConfig, option)
	if err != nil {
		return nil, err
	}

	grammarSQL.Read = read
	grammarSQL.ReadConfig = readConfig
	grammarSQL.Quoter.Bind(write, option.Prefix, read)
	return grammarSQL, nil
}

// OnConnected the event will be triggered when db server was connected
func (grammarSQL MySQL) OnConnected() error {
	version, err := grammarSQL.GetVersion()
	if err != nil {
		return err
	}
	ver577, err := semver.Make("5.7.7")
	if err != nil {
		return err
	}
	if version.LE(ver577) {
		grammarSQL.DB.Exec("SET GLOBAL innodb_file_format=`BARRACUDA`")
		grammarSQL.DB.Exec("SET GLOBAL innodb_file_per_table=`ON`;")
		grammarSQL.DB.Exec("SET GLOBAL innodb_large_prefix=`ON`;")
	}

	// Auto set sql mode
	grammarSQL.DB.Exec("SET GLOBAL sql_mode=`STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION`;")
	return nil
}

// New Create a new MySQL grammar inteface
func New(opts ...sql.Option) dbal.Grammar {
	my := MySQL{
		SQL: sql.NewSQL(&Quoter{}, opts...),
	}
	if my.Driver == "" || my.Driver == "sql" {
		my.Driver = "mysql"
	}
	// set fliptypes
	flipTypes, ok := utils.MapFilp(my.Types)
	if ok {
		my.FlipTypes = flipTypes.(map[string]string)
		my.FlipTypes["DATETIME"] = "dateTime"
		my.FlipTypes["TIME"] = "time"
		my.FlipTypes["TIMESTAMP"] = "timestamp"
	}
	return my
}
