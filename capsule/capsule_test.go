package capsule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/xun/unit"
)

func TestAdd(t *testing.T) {
	unit.SetLogger()
	m1, err := Add("test", unit.Driver(), unit.DSN())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(m1.Pool.Primary))

	m2, err := m1.Add("test2", "mysql", "root:123456@tcp(1.2.3.4:3306)/xun?charset=utf8mb4&parseTime=True&loc=Local", false)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(m2.Pool.Primary))
	assert.Equal(t, m1, m2)
}

func TestAddRead(t *testing.T) {
	unit.SetLogger()
	m1, err := AddRead("test", unit.Driver(), unit.DSN())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(m1.Pool.Readonly))

	m2, err := m1.Add("test2", "mysql", "root:123456@tcp(1.2.3.4:3306)/xun?charset=utf8mb4&parseTime=True&loc=Local", true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(m2.Pool.Readonly))
	assert.Equal(t, m1, m2)
}

func TestPing(t *testing.T) {
	unit.SetLogger()
	m1, err := Add("test", unit.Driver(), unit.DSN())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(m1.Pool.Primary))

	conn, err := m1.Primary()
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Ping(2 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := Add("test2", "mysql", "root:123456@tcp(1.2.3.4:3306)/xun?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}

	conn, err = m2.Primary()
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Ping(1 * time.Second)
	assert.Equal(t, "context deadline exceeded", err.Error())
}
