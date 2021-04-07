package query

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/yaoapp/xun/utils"
)

// DD Die and dump the current SQL and bindings.
func (builder *Builder) DD() {
	defer os.Exit(0)
	builder.Dump()
	os.Exit(0)
}

// Dump Dump the current SQL and bindings.
func (builder *Builder) Dump() {
	defer catch()
	fmt.Println(builder.ToSQL())
	utils.Println(builder.GetBindings())
	utils.Println(builder.MustGet())
}

// catch and out
func catch() {
	if r := recover(); r != nil {
		switch r.(type) {
		case string:
			color.Red("%s\n", r)
			break
		case error:
			color.Red("%s\n", r.(error).Error())
			break
		default:
			color.Red("%#v\n", r)
		}
		fmt.Println(string(debug.Stack()))
	}
}
