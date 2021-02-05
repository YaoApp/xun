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
	schema := manager.Schema()
	schema.Create()
	schema.Drop()

	qb := manager.Query()
	qb.Where()
	qb.Join()

}

func TestAddConnectionSqlite(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnectionSqlite ====================\n")
	manager := New().
		AddConnection(Config{
			Driver: "sqlite",
			DBName: "unit-test.db",
			Memory: true,
		})
	schema := manager.Schema()
	schema.Create()
	schema.Drop()
}

func TestGlobal(t *testing.T) {
	fmt.Printf("\n\n== TestGlobal ====================\n")
	schema := Schema()
	schema.Create()
	schema.Drop()

	qb := Query()
	qb.Where()
	qb.Join()
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

	schema := Schema()
	schema.Create()
	schema.Drop()

	qb := Query()
	qb.Where()
	qb.Join()
}
