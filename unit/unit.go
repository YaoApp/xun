package unit

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/yaoapp/xun/logger"
)

// "root:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local"

// Always run Always
var Always bool = true

// Is the DSN should be the name given
func Is(name string) bool {
	return os.Getenv("XUN_UNIT_NAME") == name
}

// Not the DSN should not be the  name given
func Not(name string) bool {
	return os.Getenv("XUN_UNIT_NAME") != name
}

//DriverIs the drvier should be the name given
func DriverIs(name string) bool {
	return os.Getenv("XUN_UNIT_DRIVER") == name
}

//DriverNot the drvier should not be the name given
func DriverNot(name string) bool {
	return os.Getenv("XUN_UNIT_DRIVER") != name
}

// DSN get the connection source from evn
func DSN() string {
	source := os.Getenv("XUN_UNIT_SOURCE")
	if source == "" {
		panic(errors.New("DSN does not found"))
	}
	return source
}

// Driver get the driver name from evn
func Driver() string {
	driver := os.Getenv("XUN_UNIT_DRIVER")
	if driver == "" {
		panic(errors.New("DSN does not found"))
	}
	return driver
}

// SetLogger set the unit file logger
func SetLogger() {
	logfile := os.Getenv("XUN_UNIT_LOG")
	output := os.Stdout
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			output = f
		}
	}
	logger.DefaultLogger.SetOutput(output)
	logger.DefaultErrorLogger.SetOutput(output)
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
