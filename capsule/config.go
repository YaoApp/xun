package capsule

import "fmt"

var drivers map[string]string = map[string]string{
	"mysql":  "mysql",
	"sqlite": "sqlite3",
}

// DataSource get the data source format text
func (config Config) DataSource() string {
	switch config.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		)
	case "sqlite":
		if config.Memory {
			config.Path = ":memory:"
		}
		return fmt.Sprintf("%s/%s", config.Path, config.DBName)
	default:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		)
	}
}

// DriverName get the driver name
func (config Config) DriverName() string {
	return drivers[config.Driver]
}
