package mysql

import (
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
)

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	return Mysql{
		SQL: sql.SQL{
			Driver: "mysql",
		},
	}
}
