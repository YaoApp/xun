package dbal

// Config the Connection Configuration
type Config struct {
	Driver   string `json:"driver"`        // The driver name. mysql,oci,pgsql,sqlsrv,sqlite
	DSN      string `json:"dsn,omitempty"` // The driver wrapper. sqlite:///:memory:, mysql://localhost:4486/foo?charset=UTF8
	Name     string `json:"name,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`
}

// DBConfig the database configure
type DBConfig struct {
	TablePrefix string `json:"table_prefix,omitempty"`
	DBPrefix    string `json:"db_prefix,omitempty"`
	DBName      string `json:"db,omitempty"`
	Collation   string `json:"collation,omitempty"`
	Charset     string `json:"charset,omitempty"`
}
