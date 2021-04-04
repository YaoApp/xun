package sqlite3

import (
	"fmt"

	"github.com/yaoapp/xun/grammar/sql"
)

// Quoter the database quoting query text SQL type
type Quoter struct {
	sql.Quoter
}

// WrapUnion a union subquery in parentheses.
func (quoter *Quoter) WrapUnion(sql string) string {
	return fmt.Sprintf("select * from (%s)", sql)
}
