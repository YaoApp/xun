package capsule

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

// Manager The database manager
type Manager struct {
	Pool        *Pool
	Connections *sync.Map // map[string]*Connection
}

// Pool the connection pool
type Pool struct {
	Primary  []*Connection
	Readonly []*Connection
}

// Connection The database connection
type Connection struct {
	sqlx.DB
	Config *Config
}

// Config the Connection Configuration
type Config struct {
	Driver          string     `json:"driver"`        // The driver name. mysql,oci,pgsql,sqlsrv,sqlite
	DSN             string     `json:"dsn,omitempty"` // The driver wrapper. sqlite:///:memory:, mysql://localhost:4486/foo?charset=UTF8
	Host            string     `json:"host,omitempty"`
	Port            int        `json:"port,omitempty"`
	Socket          string     `json:"socket,omitempty"`
	DBName          string     `json:"db,omitempty"`
	User            string     `json:"user,omitempty"`
	Password        string     `json:"password,omitempty"`
	Charset         string     `json:"charset,omitempty"`
	Path            string     `json:"path,omitempty"`
	ServiceName     string     `json:"service_name,omitempty"`
	InstanceName    string     `json:"instance_name,omitempty"`
	ApplicationName string     `json:"application_name,omitempty"`
	ConnectString   string     `json:"Connect_string,omitempty"`
	Service         *bool      `json:"service,omitempty"`
	Persistent      *bool      `json:"persistent,omitempty"`
	Pooled          *bool      `json:"pooled,omitempty"`
	Memory          bool       `json:"memory,omitempty"`
	SSL             *ConfigSSL `json:"ssl,omitempty"`
	ReadOnly        bool       `json:"readonly,omitempty"`
	Name            string     `json:"name,omitempty"`
}

// ConfigSSL  the Connection SSL Configuration
type ConfigSSL struct {
	Mode     string `json:"mode,omitempty"` // Determines whether or with what priority a SSL TCP/IP connection will be negotiated with the server. See the list of available modes: `https://www.postgresql.org/docs/9.4/static/libpq-connect.html#LIBPQ-CONNECT-SSLMODE`
	RootCert string `json:"root_cert,omitempty"`
	Cert     string `json:"cert,omitempty"`
	Key      string `json:"key,omitempty"`
	CAName   string `json:"ca_name,omitempty"`
	CAPath   string `json:"ca_path,omitempty"`
	Cipher   string `json:"cipher,omitempty"`
	CRL      string `json:"crl,omitempty"`
}
