package capsule

import (
	"testing"
)

func TestAddConnection(t *testing.T) {
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
}
