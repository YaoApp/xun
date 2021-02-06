package schema

import (
	"errors"
	"fmt"
	"strings"
)

// NewBlueprint create a new blueprint intance
func NewBlueprint(name string, builder *Builder) *Blueprint {
	return &Blueprint{
		Name:      name,
		Builder:   builder,
		Columns:   []*Column{},
		ColumnMap: map[string]*Column{},
		Indexes:   []*Index{},
		IndexMap:  map[string]*Index{},
	}
}

// Exists check if the table is exists
func (table *Blueprint) Exists() bool {
	row := table.validate().Builder.Conn.Write.
		QueryRowx(table.sqlExists())
	if row.Err() != nil {
		panic(row.Err())
	}

	res, err := row.SliceScan()
	if err != nil {
		return false
	}
	return table.Name == fmt.Sprintf("%s", res[0])
}

// Create a new table on the schema
func (table *Blueprint) Create() error {
	_, err := table.validate().Builder.Conn.Write.
		Exec(table.sqlCreate())
	return err
}

// MustCreate a new table on the schema
func (table *Blueprint) MustCreate() *Blueprint {
	err := table.Create()
	if err != nil {
		panic(err)
	}
	return table
}

// Drop a table from the schema.
func (table *Blueprint) Drop() error {
	_, err := table.validate().Builder.Conn.Write.
		Exec(table.sqlDrop())
	return err
}

// MustDrop a table from the schema.
func (table *Blueprint) MustDrop() {
	err := table.Drop()
	if err != nil {
		panic(err)
	}
}

// DropIfExists drop the table if the table exists
func (table *Blueprint) DropIfExists() error {
	_, err := table.validate().Builder.Conn.Write.
		Exec(table.sqlDropIfExists())
	return err
}

// MustDropIfExists drop the table if the table exists
func (table *Blueprint) MustDropIfExists() {
	err := table.DropIfExists()
	if err != nil {
		panic(err)
	}
}

// Rename a table on the schema.
func (table *Blueprint) Rename(name string) error {
	_, err := table.validate().Builder.Conn.Write.
		Exec(table.sqlRename(name))
	table.Name = name
	return err
}

// MustRename a table on the schema.
func (table *Blueprint) MustRename(name string) *Blueprint {
	err := table.Rename(name)
	if err != nil {
		panic(err)
	}
	return table
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Blueprint) BigInteger() {}

// String Create a new string column on the table.
func (table *Blueprint) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = &length
	column.Type = "string"
	table.addColumn(column)
	return column
}

// NewIndex New index instance
func (table *Blueprint) NewIndex(name string, columns ...*Column) *Index {
	index := &Index{Name: name, Columns: []*Column{}}
	index.Columns = append(index.Columns, columns...)
	index.Table = table
	return index
}

// NewColumn New column instance
func (table *Blueprint) NewColumn(name string) *Column {
	return &Column{Name: name, Table: table}
}

func (table *Blueprint) addColumn(column *Column) {
	column.validate()
	table.Columns = append(table.Columns, column)
	table.ColumnMap[column.Name] = column
}

func (table *Blueprint) addIndex(index *Index) {
	index.validate()
	table.Indexes = append(index.Table.Indexes, index)
	table.IndexMap[index.Name] = index
}

func (table *Blueprint) validate() *Blueprint {
	if table.Builder == nil {
		err := errors.New("the table " + table.Name + "does not bind the builder")
		panic(err)
	}
	return table
}

func tableNameEscaped(name string) string {
	return strings.ReplaceAll(name, "`", "")
}

func (table *Blueprint) nameEscaped() string {
	return tableNameEscaped(table.Name)
}

func (table *Blueprint) sqlExists() string {
	sql := fmt.Sprintf("SHOW TABLES like '%s'", table.nameEscaped())
	return sql
}

func (table *Blueprint) sqlDrop() string {
	sql := fmt.Sprintf("DROP TABLE `%s`", table.nameEscaped())
	return sql
}

func (table *Blueprint) sqlRename(name string) string {
	sql := fmt.Sprintf("RENAME TABLE `%s` TO `%s`", table.nameEscaped(), tableNameEscaped(name))
	return sql
}

func (table *Blueprint) sqlDropIfExists() string {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table.nameEscaped())
	return sql
}

func (table *Blueprint) sqlCreate() string {

	sql := fmt.Sprintf("CREATE TABLE `%s` (\n", table.nameEscaped())

	// columns
	stmts := []string{}
	for _, column := range table.Columns {
		stmts = append(stmts, column.sqlCreate())
	}

	for _, index := range table.Indexes {
		stmts = append(stmts, index.sqlCreate())
	}

	// indexes
	sql = sql + strings.Join(stmts, ",\n")
	sql = sql + "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

	return sql
}
