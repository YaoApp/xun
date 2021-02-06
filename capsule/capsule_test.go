package capsule

import (
	"fmt"
	"testing"
)

func TestAddConnection(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnection ====================\n")
	manager := New().
		AddConnection(Config{
			Driver:   "mysql",
			Host:     "192.168.31.119",
			Port:     3306,
			User:     "root",
			Password: "123456",
			DBName:   "xiang",
			Charset:  "utf8",
		}).
		AddConnection(Config{
			Driver:   "mysql",
			Host:     "192.168.31.119",
			Port:     3306,
			User:     "xiang",
			Password: "123456",
			DBName:   "xiang",
			Charset:  "utf8",
			ReadOnly: true,
		})

	manager.Schema()
	manager.Query()

}

func TestAddConnectionSqlite(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnectionSqlite ====================\n")
	manager := New().
		AddConnection(Config{
			Driver: "sqlite",
			DBName: "unit-test.db",
			Memory: true,
		})
	manager.Schema()
}

func TestGlobal(t *testing.T) {
	fmt.Printf("\n\n== TestGlobal ====================\n")
	Schema()
	Query()
}

func TestSetAsGlobal(t *testing.T) {
	fmt.Printf("\n\n== TestGlobalAfterSet ====================\n")
	manager := New().
		AddConnection(Config{
			Driver: "sqlite",
			DBName: "unit-test-2.db",
			Memory: true,
		})
	manager.SetAsGlobal()

	Schema()
	Query()
}
