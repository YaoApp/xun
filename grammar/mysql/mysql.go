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

// Setup the method will be executed when db server was connected
func (grammarSQL *MySQL) Setup(db *sqlx.DB, config *dbal.Config, option *dbal.Option) error {
	grammarSQL.DB = db
	grammarSQL.Config = config
	grammarSQL.Option = option
	if grammarSQL.Config == nil {
		return fmt.Errorf("config is nil")
	}
	cfg, err := mysql.ParseDSN(grammarSQL.Config.DSN)
	if err != nil {
		return err
	}
	grammarSQL.DatabaseName = cfg.DBName
	grammarSQL.SchemaName = grammarSQL.DatabaseName
	return nil
}

// OnConnected the event will be triggered when db server was connected
func (grammarSQL *MySQL) OnConnected() error {
	version, err := grammarSQL.GetVersion()
	if err != nil {
		panic(fmt.Errorf("OnConnected: %s", err))
	}
	ver577, err := semver.Make("5.7.7")
	if err != nil {
		panic(fmt.Errorf("OnConnected: %s", err))
	}
	if version.LE(ver577) {
		grammarSQL.DB.Exec("SET GLOBAL innodb_file_format=`BARRACUDA`")
		grammarSQL.DB.Exec("SET GLOBAL innodb_file_per_table=`ON`;")
		grammarSQL.DB.Exec("SET GLOBAL innodb_large_prefix=`ON`;")
	}
	return nil
}

// New Create a new MySQL grammar inteface
func New() dbal.Grammar {
	my := MySQL{
		SQL: sql.NewSQL(Quoter{}),
	}
	my.Driver = "mysql"
	my.Quoter = Quoter{}
	// set fliptypes
	flipTypes, ok := utils.MapFilp(my.Types)
	if ok {
		my.FlipTypes = flipTypes.(map[string]string)
		my.FlipTypes["DATETIME"] = "dateTime"
		my.FlipTypes["TIME"] = "time"
		my.FlipTypes["TIMESTAMP"] = "timestamp"
	}
	return &my
}
