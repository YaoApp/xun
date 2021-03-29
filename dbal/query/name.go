package query

import (
	"fmt"
	"strings"

	"github.com/yaoapp/xun/utils"
)

// TableFullName get the name name with prefix
func (name Name) TableFullName() string {
	return fmt.Sprintf("%s%s", utils.StringVal(name.Prefix), name.Name)
}

// TableAlias get the alias of name
func (name Name) TableAlias() string {
	return name.Alias
}

// Name set the Name attribute
func (builder *Builder) Name(fullname string) Name {
	name := Name{}
	name.Prefix = &builder.Conn.Option.Prefix
	namer := strings.Split(strings.ToLower(fullname), " as ")
	if len(namer) == 2 {
		name.Name = strings.Trim(namer[0], " ")
		name.Alias = strings.Trim(namer[1], " ")
		return name
	}
	name.Name = strings.Trim(fullname, " ")
	return name
}
