package unit

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/yaoapp/xun/capsule"
)

// "root:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local"

// DSN the dsns for testing
var DSN map[string]string = map[string]string{
	"mysql":  os.Getenv("XUN_UNIT_MYSQL_DSN"),
	"sqlite": os.Getenv("XUN_UNIT_SQLITE_DSN"),
	"pgsql":  os.Getenv("XUN_UNIT_POSTGRE_DSN"),
	"oracle": os.Getenv("XUN_UNIT_ORACLE_DSN"),
	"sqlsvr": os.Getenv("XUN_UNIT_SQLSVR_DSN"),
}

// Use create a capsule intance using DSN
func Use(name string) *capsule.Manager {
	dsn, has := DSN[name]
	if !has || dsn == "" {
		err := errors.New("dsn not found!" + name)
		panic(err)
	}
	return capsule.AddConnection(capsule.Config{
		Driver: name,
		DSN:    dsn,
	})
}

// Catch and out
func Catch() {
	if r := recover(); r != nil {
		switch r.(type) {
		case string:
			fmt.Printf("%s\n", r)
			break
		case error:
			fmt.Printf("%s\n", r.(error).Error())
			break
		default:
			fmt.Printf("%#v\n", r)
		}
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
	}
}
