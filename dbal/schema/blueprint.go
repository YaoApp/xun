package schema

import (
	"errors"
	"sync"

	"github.com/yaoapp/xun/grammar"
)

// NewBlueprint create a new blueprint intance
func NewBlueprint(name string, builder *Builder) *Table {
	table := &Table{
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

// Column get the column instance of the table, if the column does not exist create
func (table *Table) Column(name string) *Column {
	column, has := table.ColumnMap[name]
	if has {
		return column
	}
	return table.NewColumn(name)
}

// GetColumnListing Get the column listing for the table
func (table *Table) GetColumnListing() {
}

// GetIndexListing Get the index listing for the table
func (table *Table) GetIndexListing() {
}

// NewIndex Create a new index instance
func (table *Table) NewIndex(name string, columns ...*Column) *Index {
	index := &Index{Name: name, Columns: []*Column{}}
	index.Columns = append(index.Columns, columns...)
	index.Table = table
	return index
}

// GetIndex get the index instance for the given name, create if not exists.
func (table *Table) GetIndex(name string) *Index {
	index, has := table.IndexMap[name]
	if !has {
		index = table.NewIndex(name)
	}
	return index
}

// NewColumn Create a new column instance
func (table *Table) NewColumn(name string) *Column {
	return &Column{Name: name, Table: table}
}

// GetColumn get the column instance for the given name, create if not exists.
func (table *Table) GetColumn(name string) *Column {
	column, has := table.ColumnMap[name]
	if !has {
		column = table.NewColumn(name)
	}
	return column
}

// HasColumn Determine if the table has a given column.
func (table *Table) HasColumn(name ...string) bool {
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
func (table *Table) DropColumn(name ...string) {
	for _, n := range name {
		column := table.GetColumn(n)
		column.Drop()
	}
	table.onChange("DropColumn", name)
}

// RenameColumn Indicate that the given column should be renamed.
func (table *Table) RenameColumn(old string, new string) *Column {
	column := table.GetColumn(old)
	column.Rename(new)
	table.onChange("RenameColumn", old, new)
	return column
}

// DropIndex Indicate that the given indexes should be dropped.
func (table *Table) DropIndex(name ...string) {
	for _, n := range name {
		index := table.GetIndex(n)
		index.Drop()
	}
	table.onChange("DropIndex", name)
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Table) RenameIndex(old string, new string) *Index {
	index := table.GetIndex(old)
	index.Rename(new)
	table.onChange("RenameIndex", old, new)
	return index
}

func (table *Table) addColumn(column *Column) *Column {
	column.validate()
	table.Columns = append(table.Columns, column)
	table.ColumnMap[column.Name] = column
	table.onChange("addColumn", column)
	return column
}

func (table *Table) addIndex(index *Index) *Index {
	index.validate()
	table.Indexes = append(index.Table.Indexes, index)
	table.IndexMap[index.Name] = index
	table.onChange("addIndex", index)
	return index
}

func (table *Table) validate() *Table {
	if table.Builder == nil {
		err := errors.New("the table " + table.Name + "does not bind the builder")
		panic(err)
	}
	return table
}

// onChange call this when the table changed
func (table *Table) onChange(event string, args ...interface{}) {
}

// GrammarTable translate type to the grammar table.
func (table *Table) GrammarTable() *grammar.Table {
	db := table.Builder.Conn.WriteConfig.DBName()
	gtable := &grammar.Table{
		DBName:    db,
		Name:      table.Name,
		Comment:   table.Comment,
		Columns:   []*grammar.Column{},
		Indexes:   []*grammar.Index{},
		Collation: table.Builder.Conn.Config.Collation,
		Charset:   table.Builder.Conn.Config.Charset,
	}

	// translate columns
	ColumnMap := sync.Map{}
	for _, column := range table.Columns {
		Column := &grammar.Column{
			DBName:            db,
			TableName:         table.Name,
			Name:              column.Name,
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
		gtable.Columns = append(gtable.Columns, Column)
		ColumnMap.Store(column.Name, Column)
	}

	// translate indexes
	for _, index := range table.Indexes {
		gindex := &grammar.Index{
			DBName:    db,
			TableName: table.Name,
			Name:      index.Name,
			Type:      index.Type,
			Comment:   index.Comment,
			Columns:   []*grammar.Column{},
		}
		// bind columns and indexes
		for _, column := range index.Columns {
			Column, has := ColumnMap.Load(column.Name)
			if has {
				gindex.Columns = append(gindex.Columns, Column.(*grammar.Column))
				Column.(*grammar.Column).Indexes = append(Column.(*grammar.Column).Indexes, gindex)
			}
		}
		gtable.Indexes = append(gtable.Indexes, gindex)
	}

	gtable.SetDefaultEngine("InnoDB")
	gtable.SetDefaultCollation("utf8mb4_unicode_ci")
	gtable.SetDefaultCharset("utf8mb4")
	return gtable
}

// func (table *Table) sqlColumns() string {
// 	// SELECT *
// 	// FROM INFORMATION_SCHEMA.COLUMNS
// 	// WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME ='products';
// 	Columns := []string{
// 		"TABLE_SCHEMA AS `db_name`",
// 		"TABLE_NAME AS `table_name`",
// 		"COLUMN_NAME AS `Column`",
// 		"ORDINAL_POSITION AS `position`",
// 		"COLUMN_DEFAULT AS `default`",
// 		"IS_NULLABLE as `nullable`",
// 		"DATA_TYPE as `type`",
// 		"CHARACTER_MAXIMUM_LENGTH as `length`",
// 		"CHARACTER_OCTET_LENGTH as `octet_length`",
// 		"NUMERIC_PRECISION as `precision`",
// 		"NUMERIC_SCALE as `scale`",
// 		"DATETIME_PRECISION as `datetime_precision`",
// 		"CHARACTER_SET_NAME as `character`",
// 		"COLLATION_NAME as `collation`",
// 		"COLUMN_KEY as `key`",
// 		"EXTRA as `extra`",
// 		"COLUMN_COMMENT as `comment`",
// 	}
// 	sql := `
// 		SELECT %s
// 		FROM INFORMATION_SCHEMA.COLUMNS
// 		WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME ='%s';
// 	`
// 	cfg := table.Builder.Conn.WriteConfig
// 	fmt.Printf("sqlColumns: %#v\n", cfg.Sqlite3DBName())
// 	sql = fmt.Sprintf(sql, strings.Join(Columns, ","), "xiang", table.Name)
// 	fmt.Printf("%s", sql)
// 	return sql
// }
