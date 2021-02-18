package dbal

// Config the Connection configuration
type Config struct {
	Driver   string `json:"driver"`        // The driver name. mysql,oci,pgsql,sqlsrv,sqlite
	DSN      string `json:"dsn,omitempty"` // The driver wrapper. sqlite:///:memory:, mysql://localhost:4486/foo?charset=UTF8
	Name     string `json:"name,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`
}

// Option the database configuration
type Option struct {
	Prefix    string `json:"prefix,omitempty"` // Table prifix
	Collation string `json:"collation,omitempty"`
	Charset   string `json:"charset,omitempty"`
}
