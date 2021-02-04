package capsule

import "fmt"

// DataSource get the data source format text
func (config Config) DataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}

// DriverName get the driver name
func (config Config) DriverName() string {
	return config.Driver
}
