package sql

import (
	"fmt"

	"github.com/yaoapp/xun/grammar"
)

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	return SQL{
		Driver: "sql",
	}
}

// Exists the Exists
func (grammar SQL) Exists(table *grammar.Table) string {
	sql := fmt.Sprintf("SHOW TABLES like '%s'", table.Name)
	return sql
}
