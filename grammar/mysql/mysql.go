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
	dbal.Register("mysql", New(), dbal.Hook{
		OnConnected: func(grammarSQL dbal.Grammar, db *sqlx.DB) {
			version, err := grammarSQL.GetVersion(db)
			if err != nil {
				panic(fmt.Errorf("OnConnected: %s", err))
			}
			ver577, err := semver.Make("5.7.7")
			if err != nil {
				panic(fmt.Errorf("OnConnected: %s", err))
			}
			if version.LE(ver577) {
				db.Exec("SET GLOBAL innodb_file_format=`BARRACUDA`")
				db.Exec("SET GLOBAL innodb_file_per_table=`ON`;")
				db.Exec("SET GLOBAL innodb_large_prefix=`ON`;")
			}
		},
	})
}

// Config set the configure using DSN
func (grammarSQL *MySQL) Config(dsn string) {
	grammarSQL.DSN = dsn
	cfg, err := mysql.ParseDSN(grammarSQL.DSN)
	if err != nil {
		panic(err)
	}
	grammarSQL.DB = cfg.DBName
	grammarSQL.Schema = grammarSQL.DB
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
