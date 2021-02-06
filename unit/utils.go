package unit

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/yaoapp/xun/capsule"
)

// "root:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local"

// DSN the dsns for testing
var dsns map[string]string = map[string]string{
	"mysql":  os.Getenv("XUN_UNIT_MYSQL_DSN"),
	"sqlite": os.Getenv("XUN_UNIT_SQLITE_DSN"),
	"pgsql":  os.Getenv("XUN_UNIT_POSTGRE_DSN"),
	"oracle": os.Getenv("XUN_UNIT_ORACLE_DSN"),
	"sqlsvr": os.Getenv("XUN_UNIT_SQLSVR_DSN"),
}

// Use create a capsule intance using DSN
func Use(name string) *capsule.Manager {
	dsn := DSN(name)
	return capsule.AddConnection(capsule.Config{
		Driver: name,
		DSN:    dsn,
	})
}

// DSN get the dsn from evn
func DSN(name string) string {
	dsn, has := dsns[name]
	if !has || dsn == "" {
		err := errors.New("dsn not found!" + name)
		panic(err)
	}
	return dsn
}

// Catch and out
func Catch() {
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
