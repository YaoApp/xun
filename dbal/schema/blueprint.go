package schema

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/yaoapp/xun/grammar"
)

// NewBlueprint create a new blueprint intance
func NewBlueprint(name string, builder *Builder) *Blueprint {
	table := &Blueprint{
		Name:      name,
		Builder:   builder,
		Columns:   []*Column{},
		ColumnMap: map[string]*Column{},
		Indexes:   []*Index{},
		IndexMap:  map[string]*Index{},
	}
	table.onChange("NewBlueprint", name, builder)
	return table
}

// Alter a table one the schema
func (table *Blueprint) Alter(callback func(table *Blueprint)) error {
	alter := table.Builder.Table(table.Name)
	alter.alter = true
	callback(alter)
	alter.sqlAlter()
	return nil
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

// Column get the column instance of the table, if the column does not exist create
func (table *Blueprint) Column(name string) *Column {
	column, has := table.ColumnMap[name]
	if has {
		return column
	}
	return table.NewColumn(name)
}

// GetColumnListing Get the column listing for the table
func (table *Blueprint) GetColumnListing() {
	rows, err := table.validate().Builder.Conn.Write.
		Queryx(table.sqlColumns())
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		row := &TableField{}
		err = rows.StructScan(row)
		if err != nil {
			panic(err)
		}
		if table.HasColumn(row.Field) {
			table.Column(row.Field).UpField(row)
		} else {
			column := table.NewColumn(row.Field)
			column.UpField(row)
			table.addColumn(column)
		}
	}
}

// GetIndexListing Get the index listing for the table
func (table *Blueprint) GetIndexListing() {
}

// NewIndex Create a new index instance
func (table *Blueprint) NewIndex(name string, columns ...*Column) *Index {
	index := &Index{Name: name, Columns: []*Column{}}
	index.Columns = append(index.Columns, columns...)
	index.Table = table
	return index
}

// GetIndex get the index instance for the given name, create if not exists.
func (table *Blueprint) GetIndex(name string) *Index {
	index, has := table.IndexMap[name]
	if !has {
		index = table.NewIndex(name)
	}
	return index
}

// NewColumn Create a new column instance
func (table *Blueprint) NewColumn(name string) *Column {
	return &Column{Name: name, Table: table}
}

// GetColumn get the column instance for the given name, create if not exists.
func (table *Blueprint) GetColumn(name string) *Column {
	column, has := table.ColumnMap[name]
	if !has {
		column = table.NewColumn(name)
	}
	return column
}

// HasColumn Determine if the table has a given column.
func (table *Blueprint) HasColumn(name ...string) bool {
	has := true
	for _, n := range name {
		_, has = table.ColumnMap[n]
		if !has {
			return has
		}
	}
	return has
}

// DropColumn Indicate that the given columns should be dropped.
func (table *Blueprint) DropColumn(name ...string) {
	for _, n := range name {
		column := table.GetColumn(n)
		column.Drop()
	}
	table.onChange("DropColumn", name)
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Blueprint) RenameColumn(old string, new string) *Column {
	column := table.GetColumn(old)
	column.Rename(new)
	table.onChange("RenameColumn", old, new)
	return column
}

// DropIndex Indicate that the given indexes should be dropped.
func (table *Blueprint) DropIndex(name ...string) {
	for _, n := range name {
		index := table.GetIndex(n)
		index.Drop()
	}
	table.onChange("DropIndex", name)
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Blueprint) RenameIndex(old string, new string) *Index {
	index := table.GetIndex(old)
	index.Rename(new)
	table.onChange("RenameIndex", old, new)
	return index
}

func (table *Blueprint) addColumn(column *Column) *Column {
	column.validate()
	table.Columns = append(table.Columns, column)
	table.ColumnMap[column.Name] = column
	table.onChange("addColumn", column)
	return column
}

func (table *Blueprint) addIndex(index *Index) *Index {
	index.validate()
	table.Indexes = append(index.Table.Indexes, index)
	table.IndexMap[index.Name] = index
	table.onChange("addIndex", index)
	return index
}

func (table *Blueprint) validate() *Blueprint {
	if table.Builder == nil {
		err := errors.New("the table " + table.Name + "does not bind the builder")
		panic(err)
	}
	return table
}

// onChange call this when the table changed
func (table *Blueprint) onChange(event string, args ...interface{}) {
}

// GrammarTable translate type to the grammar table.
func (table *Blueprint) GrammarTable() *grammar.Table {
	db := table.Builder.Conn.WriteConfig.DBName()
	gtable := &grammar.Table{
		DBName:    db,
		Name:      table.Name,
		Comment:   table.Comment,
		Fields:    []*grammar.Field{},
		Indexes:   []*grammar.Index{},
		Collation: table.Builder.Conn.Config.Collation,
		Charset:   table.Builder.Conn.Config.Charset,
	}

	// translate columns
	fieldMap := sync.Map{}
	for _, column := range table.Columns {
		field := &grammar.Field{
			DBName:            db,
			TableName:         table.Name,
			Field:             column.Name,
			Type:              column.Type,
			Length:            column.Length,
			Precision:         column.Precision(),
			Scale:             column.Scale(),
			DatetimePrecision: column.DatetimePrecision(),
			Comment:           column.Comment,
			Collation:         table.Builder.Conn.Config.Collation,
			Charset:           table.Builder.Conn.Config.Charset,
			Indexes:           []*grammar.Index{},
		}
		gtable.Fields = append(gtable.Fields, field)
		fieldMap.Store(column.Name, field)
	}

	// translate indexes
	for _, index := range table.Indexes {
		gindex := &grammar.Index{
			DBName:    db,
			TableName: table.Name,
			Index:     index.Name,
			Type:      index.Type,
			Comment:   index.Comment,
			Fields:    []*grammar.Field{},
		}
		// bind columns and indexes
		for _, column := range index.Columns {
			field, has := fieldMap.Load(column.Name)
			if has {
				gindex.Fields = append(gindex.Fields, field.(*grammar.Field))
				field.(*grammar.Field).Indexes = append(field.(*grammar.Field).Indexes, gindex)
			}
		}
		gtable.Indexes = append(gtable.Indexes, gindex)
	}

	gtable.SetDefaultEngine("InnoDB")
	gtable.SetDefaultCollation("utf8mb4_unicode_ci")
	gtable.SetDefaultCharset("utf8mb4")
	return gtable
}

func tableNameEscaped(name string) string {
	return strings.ReplaceAll(name, "`", "")
}

func (table *Blueprint) nameEscaped() string {
	return tableNameEscaped(table.Name)
}

func (table *Blueprint) sqlColumns() string {
	// SELECT *
	// FROM INFORMATION_SCHEMA.COLUMNS
	// WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME ='products';
	fields := []string{
		"TABLE_SCHEMA AS `db_name`",
		"TABLE_NAME AS `table_name`",
		"COLUMN_NAME AS `field`",
		"ORDINAL_POSITION AS `position`",
		"COLUMN_DEFAULT AS `default`",
		"IS_NULLABLE as `nullable`",
		"DATA_TYPE as `type`",
		"CHARACTER_MAXIMUM_LENGTH as `length`",
		"CHARACTER_OCTET_LENGTH as `octet_length`",
		"NUMERIC_PRECISION as `precision`",
		"NUMERIC_SCALE as `scale`",
		"DATETIME_PRECISION as `datetime_precision`",
		"CHARACTER_SET_NAME as `character`",
		"COLLATION_NAME as `collation`",
		"COLUMN_KEY as `key`",
		"EXTRA as `extra`",
		"COLUMN_COMMENT as `comment`",
	}
	sql := ` 
		SELECT %s	
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME ='%s';
	`
	cfg := table.Builder.Conn.WriteConfig
	fmt.Printf("sqlColumns: %#v\n", cfg.Sqlite3DBName())
	sql = fmt.Sprintf(sql, strings.Join(fields, ","), "xiang", table.nameEscaped())
	fmt.Printf("%s", sql)
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

func (table *Blueprint) sqlAlter() []string {
	table.sqlColumns()
	return []string{}
}
