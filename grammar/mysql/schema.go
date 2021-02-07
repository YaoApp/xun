package mysql

import (
	"fmt"

	"github.com/yaoapp/xun/grammar"
)

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	return Grammar{
		Driver: "mysql",
	}
}

// Exists the Exists
func (grammar Grammar) Exists(table grammar.Table) string {
	fmt.Printf("⚠️Grammar: mysql\n")
	sql := fmt.Sprintf("SHOW TABLES like '%s'", table.Name)
	return sql
}
