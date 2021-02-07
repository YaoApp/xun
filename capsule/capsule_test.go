package capsule

import (
	"fmt"
	"testing"

	"github.com/yaoapp/xun/dbal/schema"
)

func TestAddConnection(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnection ====================\n")
	manager := New().
		AddConnection(schema.ConnConfig{
			Driver:   "mysql",
			Host:     "192.168.31.119",
			Port:     3306,
			User:     "root",
			Password: "123456",
			DBName:   "xiang",
			Charset:  "utf8",
		}).
		AddConnection(schema.ConnConfig{
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
		AddConnection(schema.ConnConfig{
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
		AddConnection(schema.ConnConfig{
			Driver: "sqlite",
			DBName: "unit-test-2.db",
			Memory: true,
		})
	manager.SetAsGlobal()

	Schema()
	Query()
}
