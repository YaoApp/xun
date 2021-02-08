package sqlite3

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/grammar/sql"
)

// New Create a new mysql grammar inteface
func New() grammar.Grammar {
	sqlite := SQLite3{
		SQL: sql.NewSQL(),
	}
	sqlite.Driver = "sqlite3"
	return sqlite
}

// Exists the Exists
func (grammar SQLite3) Exists(table string, db *sqlx.DB) bool {
	fmt.Printf("⚠️Grammar: %s\n", grammar.Driver)
	// sql := fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE type='table' AND name='%s'", table.Name)
	return false

}
