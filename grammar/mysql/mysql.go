package mysql

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //Load mysql driver
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
