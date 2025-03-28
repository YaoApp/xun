package dbal

import (
	"time"

	"github.com/blang/semver/v4"
	"github.com/jmoiron/sqlx"
)

// Connection the database connection
type Connection struct {
	DB      *sqlx.DB
	Config  *Config
	Option  *Option
	Version *Version
}

// Config the Connection configuration
type Config struct {
	Driver   string `json:"driver"`        // The driver name. mysql,pgsql,sqlite3,oci,sqlsrv
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

// Version the database version
type Version struct {
	Driver string
	semver.Version
}

// Table the table struct
type Table struct {
	DBName        string    `db:"db_name"`
	SchemaName    string    `db:"schema_name"`
	TableName     string    `db:"table_name"`
	Comment       string    `db:"table_comment"`
	Type          string    `db:"table_type"`
	Engine        string    `db:"engine"`
	CreateTime    time.Time `db:"create_time"`
	CreateOptions string    `db:"create_options"`
	Collation     string    `db:"collation"`
	Charset       string    `db:"charset"`
	Rows          int       `db:"table_rows"`
	RowLength     int       `db:"avg_row_length"`
	IndexLength   int       `db:"index_length"`
	AutoIncrement int       `db:"auto_increment"`
	Primary       *Primary
	ColumnMap     map[string]*Column
	IndexMap      map[string]*Index
	Columns       []*Column
	Indexes       []*Index
	Commands      []*Command
}

// Column the table Column
type Column struct {
	DBName                   string      `db:"db_name"`
	TableName                string      `db:"table_name"`
	Name                     string      `db:"name"`
	Position                 int         `db:"position"`
	Default                  interface{} `db:"default"`
	DefaultRaw               string      `db:"default_raw"`
	Nullable                 bool        `db:"nullable"`
	IsUnsigned               bool        `db:"unsigned"`
	Type                     string      `db:"type"`
	Length                   *int        `db:"length"`
	OctetLength              *int        `db:"octet_length"`
	Precision                *int        `db:"precision"`
	Scale                    *int        `db:"scale"`
	DateTimePrecision        *int        `db:"datetime_precision"`
	Charset                  *string     `db:"charset"`
	Collation                *string     `db:"collation"`
	Key                      *string     `db:"key"`
	Extra                    *string     `db:"extra"`
	Comment                  *string     `db:"comment"`
	Primary                  bool        `db:"primary"`
	TypeName                 string      `db:"type_name"`
	MaxLength                int
	DefaultLength            int
	MaxPrecision             int
	DefaultPrecision         int
	MaxScale                 int
	DefaultScale             int
	MaxDateTimePrecision     int
	DefaultDateTimePrecision int
	Option                   []string
	Table                    *Table
	Indexes                  []*Index
	Constraint               *Constraint
}

// Index the talbe index
type Index struct {
	DBName       string  `db:"db_name"`
	TableName    string  `db:"table_name"`
	ColumnName   string  `db:"column_name"`
	Name         string  `db:"index_name"`
	SEQ          int     `db:"seq_in_index"`
	SeqColumn    int     `db:"seq_in_column"`
	Collation    *string `db:"collation"`
	Nullable     bool    `db:"nullable"`
	Unique       bool    `db:"unique"`
	Primary      bool    `db:"primary"`
	SubPart      int     `db:"sub_part"`
	Type         string  `db:"type"`
	IndexType    string  `db:"index_type"`
	Comment      *string `db:"comment"`
	IndexComment *string `db:"index_comment"`
	Table        *Table
	Columns      []*Column
}

// Primary the table primary key
type Primary struct {
	DBName    string `db:"db_name"`
	TableName string `db:"table_name"`
	Name      string `db:"primary_name"`
	Table     *Table
	Columns   []*Column
}

// Constraint the table constraint
type Constraint struct {
	SchemaName string
	TableName  string
	ColumnName string
	Type       string
	Args       []string
	Table      *Table
}

// Command The Command that should be run for the table.
type Command struct {
	Name    string        // The command name
	Params  []interface{} // The command parameters
	Success func()        // Success callback function
	Fail    func()        // Fail callback function
}

// Name the from attribute ( table_name as t1,  column_name as c1...)
type Name struct {
	Prefix string
	Name   string
	Alias  string
}

// Expression the raw query expression.
type Expression struct {
	Value interface{}
}

// Where The where constraint for the query.
type Where struct {
	Type     string // basic, nested, sub, null, notnull ...
	Column   interface{}
	First    interface{} // for comparing the two columns
	Second   interface{} // for comparing the two columns
	SQL      string
	Operator string
	Boolean  string
	Wheres   []Where
	Query    *Query
	Value    interface{}
	Values   []interface{}
	ValuesIn interface{}
	Not      bool
	Offset   int
}

// Join the join clause for the query
type Join struct {
	Type   string      // inner, left, right, cross
	Name   interface{} // The table the join clause is joining to.
	Query  *Query
	Alias  string
	SQL    interface{}
	Offset int
}

// Union the query union statement
type Union struct {
	All   bool // Union all
	Query *Query
}

// Aggregate An aggregate function and column to be run.
type Aggregate struct {
	Func    string        // AVG, COUNT, MIN, MAX, SUM
	Columns []interface{} // The columns for Aggregate
}

// Having The having constraints for the query.
type Having struct {
	Type     string
	Column   interface{}
	Operator string
	Value    interface{}
	Boolean  string
	Offset   int
	Values   []interface{}
	Not      bool
	SQL      string
}

// Order the order query
type Order struct {
	Type      string
	Column    interface{}
	Direction string
	Offset    int
	SQL       string
}

// From the from query
type From struct {
	Type   string
	Name   interface{}
	Alias  string
	Offset int
	SQL    string
}

// Select the from query
type Select struct {
	Type   string
	Name   interface{}
	Alias  string
	Offset int
	SQL    string
}

// Query the query builder
type Query struct {
	UseWriteConnection bool                     // Whether to use write connection for the select. default is false
	Lock               interface{}              //  Indicates whether row locking is being used.
	From               From                     // The table which the query is targeting.
	Columns            []interface{}            // The columns that should be returned. (Name or Expression)
	Aggregate          Aggregate                // An aggregate function and column to be run.
	Wheres             []Where                  // The where constraints for the query.
	Joins              []Join                   // The table joins for the query.
	Unions             []Union                  // The query union statements.
	UnionLimit         int                      // The maximum number of union records to return.
	UnionOffset        int                      // The number of union records to skip.
	UnionOrders        []Order                  // The orderings for the union query.
	Orders             []Order                  // The orderings for the query.
	Limit              int                      // The maximum number of records to return.
	Offset             int                      // The number of records to skip.
	Groups             []interface{}            // The groupings for the query.
	Havings            []Having                 // The having constraints for the query.
	Bindings           map[string][]interface{} // The current query value bindings.
	Distinct           bool                     // Indicates if the query returns distinct results. Occasionally contains the columns that should be distinct. default is false
	DistinctColumns    []interface{}            // Indicates if the query returns distinct results. Occasionally contains the columns that should be distinct.
	IsJoinClause       bool                     // Determine if the query is a join clause.
	BindingOffset      int                      // The Binding offset before select
	SQL                string                   // The SQL STMT
}
