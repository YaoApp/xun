package schema

import "github.com/jmoiron/sqlx"

// Connection DB Connection
type Connection struct {
	Write       *sqlx.DB
	WriteConfig *ConnConfig
	Config      *Config
}

// Schema The database Schema interface
type Schema interface {
	Table(string) *Blueprint
	HasTable(string) bool
	Create(string, func(table *Blueprint)) error
	MustCreate(string, func(table *Blueprint)) *Blueprint
	Drop(string) error
	MustDrop(string)
	DropIfExists(string) error
	MustDropIfExists(string)
	Rename(string, string) error
	MustRename(string, string) *Blueprint
	Alter(string, func(table *Blueprint)) error
	GetColumnType(string) string
	GetIndexType(string) string
}

// BlueprintMethods  the bluprint interface
type BlueprintMethods interface {
	BigInteger()
	String(name string, length int) *Blueprint
	Primary()
}

// Builder the dbal schema driver
type Builder struct {
	Conn *Connection
	Schema
}

// Config the database configure
type Config struct {
	TablePrefix string `json:"table_prefix,omitempty"`
	DBPrefix    string `json:"db_prefix,omitempty"`
	DBName      string `json:"db,omitempty"`
	Collation   string `json:"collation,omitempty"`
	Charset     string `json:"charset,omitempty"`
}

// ConnConfig the Connection Configuration
type ConnConfig struct {
	Driver          string         `json:"driver"`        // The driver name. mysql,oci,pgsql,sqlsrv,sqlite
	DSN             string         `json:"dsn,omitempty"` // The driver wrapper. sqlite:///:memory:, mysql://localhost:4486/foo?charset=UTF8
	Host            string         `json:"host,omitempty"`
	Port            int            `json:"port,omitempty"`
	Socket          string         `json:"socket,omitempty"`
	DBName          string         `json:"db,omitempty"`
	User            string         `json:"user,omitempty"`
	Password        string         `json:"password,omitempty"`
	Charset         string         `json:"charset,omitempty"`
	Path            string         `json:"path,omitempty"`
	ServiceName     string         `json:"service_name,omitempty"`
	InstanceName    string         `json:"instance_name,omitempty"`
	ApplicationName string         `json:"application_name,omitempty"`
	ConnectString   string         `json:"Connect_string,omitempty"`
	Service         *bool          `json:"service,omitempty"`
	Persistent      *bool          `json:"persistent,omitempty"`
	Pooled          *bool          `json:"pooled,omitempty"`
	Memory          bool           `json:"memory,omitempty"`
	SSL             *ConnConfigSSL `json:"ssl,omitempty"`
	ReadOnly        bool           `json:"readonly,omitempty"`
	Name            string         `json:"name,omitempty"`
}

// ConnConfigSSL  the Connection SSL Configuration
type ConnConfigSSL struct {
	Mode     string `json:"mode,omitempty"` // Determines whether or with what priority a SSL TCP/IP connection will be negotiated with the server. See the list of available modes: `https://www.postgresql.org/docs/9.4/static/libpq-connect.html#LIBPQ-CONNECT-SSLMODE`
	RootCert string `json:"root_cert,omitempty"`
	Cert     string `json:"cert,omitempty"`
	Key      string `json:"key,omitempty"`
	CAName   string `json:"ca_name,omitempty"`
	CAPath   string `json:"ca_path,omitempty"`
	Cipher   string `json:"cipher,omitempty"`
	CRL      string `json:"crl,omitempty"`
}

// Blueprint the table blueprint
type Blueprint struct {
	BlueprintMethods
	Builder   *Builder
	Comment   string
	Name      string
	Columns   []*Column
	ColumnMap map[string]*Column
	Indexes   []*Index
	IndexMap  map[string]*Index
	alter     bool
}

// Column the table column definition
type Column struct {
	Comment  string
	Name     string
	Type     string
	Length   *int
	Args     interface{}
	Default  interface{}
	Nullable *bool
	Unsigned *bool
	Table    *Blueprint
	dropped  bool
	renamed  bool
	newname  string
}

// Index  the table index definition
type Index struct {
	Comment string
	Name    string
	Type    string
	Columns []*Column
	Table   *Blueprint
	dropped bool
	renamed bool
	newname string
}

// TableField the table field
type TableField struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   interface{} `db:"Extra"`
}

// TableIndex the table index
type TableIndex struct {
	NonUnique  int
	KeyName    string
	SeqInIndex int
	ColumnName string
}
