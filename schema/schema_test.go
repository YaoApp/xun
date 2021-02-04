package schema

import (
	"testing"
)

func TestNewMySQL(t *testing.T) {
	schema := NewMySQL()
	schema.Create()
	schema.Drop()
}

func TestNewSQLite(t *testing.T) {
	schema := NewSQLite()
	schema.Create()
	schema.Drop()
}

func TestNewSQLServer(t *testing.T) {
	schema := NewSQLServer()
	schema.Create()
	schema.Drop()
}

func TestNewOracle(t *testing.T) {
	schema := NewOracle()
	schema.Create()
	schema.Drop()
}

func TestNewPostgreSQL(t *testing.T) {
	schema := NewPostgreSQL()
	schema.Create()
	schema.Drop()
}
