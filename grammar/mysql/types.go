package mysql

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
	"github.com/yaoapp/xun/utils"
)

// MySQL the MySQL Grammar
type MySQL struct {
	sql.SQL
}

// New Create a new MySQL grammar inteface
func New(dsn string) grammar.Grammar {
	my := MySQL{
		SQL: sql.NewSQL(dsn, Quoter{}),
	}
	my.Driver = "mysql"
	my.Quoter = Quoter{}

	// set fliptypes
	flipTypes, ok := utils.MapFilp(my.Types)
	if ok {
		my.FlipTypes = flipTypes.(map[string]string)
	}

	my.DBName()
	my.SchemaName()
	return &my
}
