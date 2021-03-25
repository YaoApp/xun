package query

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/utils"
)

// TableFullName get the table name with prefix
func (table Table) TableFullName() string {
	return fmt.Sprintf("%s%s", utils.StringVal(table.Prefix), table.Name)
}

// TableAlias get the alias of table
func (table Table) TableAlias() string {
	return table.Alias
}

// NewTable set the From attribute
func (builder *Builder) NewTable(name string) Table {
	table := Table{}
	table.Prefix = &builder.Conn.Option.Prefix
	namer := strings.Split(strings.ToLower(name), " as ")
	if len(namer) == 2 {
		table.Name = strings.Trim(namer[0], " ")
		table.Alias = strings.Trim(namer[1], " ")
		return table
	}
	table.Name = strings.Trim(name, " ")
	return table
}
