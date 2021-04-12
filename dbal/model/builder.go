package model

import (
	"fmt"

	"github.com/yaoapp/xun/dbal/query"
)

// Make make a new xun model instance
func Make(query query.Query, v interface{}) *Model {
	if v != nil {
		fmt.Println(v)
	}
	return nil
}
