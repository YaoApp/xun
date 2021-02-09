package grammar

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Grammar the Grammar inteface
type Grammar interface {
	Exists(name string, db *sqlx.DB) bool
	Create(table *Table, db *sqlx.DB) error
	Drop(name string, db *sqlx.DB) error
	DropIfExists(name string, db *sqlx.DB) error
	Rename(old string, new string, db *sqlx.DB) error
	// GetColumnListing(table *Table, db *sqlx.DB) []*Field
}

// Quoter the database quoting query text intrface
type Quoter interface {
	ID(name string, db *sqlx.DB) string
	VAL(v interface{}, db *sqlx.DB) string // operates on both string and []byte and int or other types.
}

// SQLBuilder the database sql gender intrface
type SQLBuilder interface {
	SQLCreateColumn(db *sqlx.DB, field *Field, types map[string]string, quoter Quoter) string
	SQLCreateIndex(db *sqlx.DB, index *Index, indexTypes map[string]string, quoter Quoter) string
	SQLTableExists(db *sqlx.DB, name string, quoter Quoter) string
}

// Table the table struct
type Table struct {
	DBName        string    `db:"db_name"`
	Name          string    `db:"table_name"`
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
	Fields        []*Field
	Indexes       []*Index
	IndexesUnique []*Index
}

// Field the table field
type Field struct {
	DBName            string      `db:"db_name"`
	TableName         string      `db:"table_name"`
	Field             string      `db:"field"`
	Position          int         `db:"position"`
	Default           interface{} `db:"default"`
	Nullable          bool        `db:"nullable"`
	Type              string      `db:"type"`
	Length            int         `db:"length"`
	OctetLength       string      `db:"octet_length"`
	Precision         int         `db:"precision"`
	Scale             int         `db:"scale"`
	DatetimePrecision int         `db:"datetime_precision"`
	Charset           string      `db:"charset"`
	Collation         string      `db:"collation"`
	Key               string      `db:"key"`
	Extra             string      `db:"extra"`
	Comment           string      `db:"comment"`
	Primary           bool        `db:"primary"`
	Indexes           []*Index
}

// Index the talbe index
type Index struct {
	DBName       string `db:"db_name"`
	TableName    string `db:"table_name"`
	Index        string `db:"index_name"`
	SEQ          int    `db:"SEQ"`
	Field        string `db:"field"`
	Collation    string `db:"collation"`
	Nullable     bool   `db:"nullable"`
	Unique       bool   `db:"unique"`
	SubPart      int    `db:"sub_part"`
	Type         string `db:"type"`
	Comment      string `db:"comment"`
	IndexComment string `db:"index_comment"`
	Fields       []*Field
}
