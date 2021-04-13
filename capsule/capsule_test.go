package capsule

import (
	"fmt"
	"testing"
)

func TestAddConnection(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnection ====================\n")
	manager := New().
		AddConn("primary", "mysql", "root:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local").
		AddReadConn("secondary", "mysql", "xiang:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local")
	manager.Schema()
	manager.Query()
	manager.Make("{}")

}

func TestAddConnectionSqlite(t *testing.T) {
	fmt.Printf("\n\n== TestAddConnectionSqlite ====================\n")
	manager := New().
		AddConn("primary", "sqlite3", "sqlite3://:memory:/capsule-test.db")
	manager.Schema()
	manager.Query()
}

func TestGlobal(t *testing.T) {
	fmt.Printf("\n\n== TestGlobal ====================\n")
	Schema()
	Query()
	Make("{}")
}

func TestSetAsGlobal(t *testing.T) {
	fmt.Printf("\n\n== TestGlobalAfterSet ====================\n")
	manager := New().
		AddConn("primary", "sqlite3", "sqlite3://:memory:/capsule-test-2.db")
	manager.SetAsGlobal()
	Schema()
	Query()
	Make("{}")
}
