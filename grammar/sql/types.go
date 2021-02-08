package sql

import "github.com/yaoapp/xun/grammar"

// SQL the SQL Grammar
type SQL struct {
	Driver string
	grammar.Grammar
}
