package hooks

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/yaoapp/xun/dbal"
	gmysql "github.com/yaoapp/xun/grammar/mysql"
	gpostgres "github.com/yaoapp/xun/grammar/postgres"
	gsql "github.com/yaoapp/xun/grammar/sql"
	gsqlite3 "github.com/yaoapp/xun/grammar/sqlite3"
)

var (
	// Hooks is the `sqlhooks.Hooks` registry
	Hooks = map[string]sqlhooks.Hooks{}

	// Drivers is the Driver registry
	Drivers = map[string]*Driver{}
)

// NoHooksError no hooks error
var NoHooksError = fmt.Errorf("no hooks error")

// RegisterHook register a `sqlhooks.Hooks`
func RegisterHook(name string, hook sqlhooks.Hooks) error {
	if _, ok := Hooks[name]; ok {
		return fmt.Errorf("hook %s already registered", name)
	}
	Hooks[name] = hook
	return nil
}

// RegisterDriver register a driver with the given name
// driver name format: type:hook1[:hook2][:...], type must be one of [mysql, postgres, sqlite3]
// it will both register the `database/sql/driver` and `yaoapp/xun/dbal`
func RegisterDriver(name string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	if _, ok := Drivers[name]; ok {
		return nil
	}
	d, err := NewDriver(name)
	if err != nil {
		if errors.Is(err, NoHooksError) {
			return nil // NOTE: no need to register if no hooks!
		}
		return err
	}
	Drivers[name] = d
	sql.Register(name, d)
	dbal.Register(name, d.grammar())
	return nil
}

// NewDriver create a new driver with `sqlhooks.Hooks` support
// Usage:
// 1. register hook, e.g. `hooks.RegisterHook("log", Default)` in `yaoapp/xun/grammar/hooks/log package`
// 2. register driver, e.g. `hooks.RegisterDriver("mysql:log")`
func NewDriver(name string) (*Driver, error) {
	parts := strings.Split(name, ":")
	typ := parts[0]
	if typ != "mysql" && typ != "postgres" && typ != "sqlite3" {
		return nil, fmt.Errorf("driver type %s not in [mysql, postgres, sqlite3]", typ)
	}
	if len(parts) == 1 {
		return nil, NoHooksError
	}
	d := &Driver{
		name: name,
		typ:  typ,
	}
	for i, hookName := range parts {
		if i == 0 { // NOTE: 0 is the default type driver, no need to repeat
			continue
		}
		if hook, ok := Hooks[hookName]; ok {
			d.Driver = sqlhooks.Wrap(d.driver(), hook)
		} else {
			return nil, fmt.Errorf("sql hook %s not found", hookName)
		}
	}
	if d.Driver == nil {
		return nil, fmt.Errorf("no driver hooks")
	}
	return d, nil
}

// Driver is the database driver which supports the `sqlhooks` in specific naming rule
type Driver struct {
	name string
	typ  string // NOTE: "oneof=mysql postgres sqlite3"`
	driver.Driver
}

// Name returns the name of the driver
func (d *Driver) Name() string {
	return d.name
}

// Type returns the type of the driver
func (d *Driver) Type() string {
	return d.typ
}

func (d *Driver) driver() driver.Driver {
	if d.Driver != nil {
		return d.Driver
	}
	switch d.typ {
	case "mysql":
		return &mysql.MySQLDriver{}
	case "postgres":
		return &pq.Driver{}
	case "sqlite3":
		return &sqlite3.SQLiteDriver{}
	}
	return nil
}

func (d *Driver) grammar() dbal.Grammar {
	switch d.typ {
	case "mysql":
		return gmysql.New(gsql.WithDriver(d.name))
	case "postgres":
		return gpostgres.New(gsql.WithDriver(d.name))
	case "sqlite3":
		return gsqlite3.New(gsql.WithDriver(d.name))
	}
	return nil
}

var _ driver.Driver = (*Driver)(nil)
