package sqlite3

import (
	"fmt"

	"github.com/yaoapp/xun/grammar"
)

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	return Grammar{
		Driver: "sqlite3",
	}
}

// Exists the Exists
func (grammar Grammar) Exists(table grammar.Table) string {
	fmt.Printf("⚠️Grammar: sqlite3\n")
	sql := fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name='%s'", table.Name)
	return sql

}
