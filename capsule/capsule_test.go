package capsule

import (
	"testing"

	"github.com/yaoapp/xun/unit"
)

func TestAddConnection(t *testing.T) {
	unit.SetLogger()
	manager := New().
		AddConn("primary", "mysql", "root:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local").
		AddReadConn("secondary", "mysql", "xiang:123456@tcp(192.168.31.119:3306)/xiang?charset=utf8mb4&parseTime=True&loc=Local")
	manager.Schema()
	manager.Query()
	manager.Make("hello", []byte("{}"))

}

func TestAddConnectionSqlite(t *testing.T) {
	manager := New().
		AddConn("primary", "sqlite3", "sqlite3://:memory:/capsule-test.db")
	manager.Schema()
	manager.Query()
}

func TestGlobal(t *testing.T) {
	Schema()
	Query()
	Make("hello", []byte("{}"))
}

func TestSetAsGlobal(t *testing.T) {
	manager := New().
		AddConn("primary", "sqlite3", "sqlite3://:memory:/capsule-test-2.db")
	manager.SetAsGlobal()
	Schema()
	Query()
	Make("hello", []byte("{}"))
}
