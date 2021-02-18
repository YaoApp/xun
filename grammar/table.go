package grammar

// NewTable create a grammar table
func NewTable(name string, schemaName string, dbName string) Table {
	return Table{
		DBName:     dbName,
		SchemaName: schemaName,
		Name:       name,
		Primary:    nil,
		Columns:    []*Column{},
		ColumnMap:  map[string]*Column{},
		Indexes:    []*Index{},
		IndexMap:   map[string]*Index{},
		Commands:   []*Command{},
	}
}

// NewIndex create a new index intstance
func (table *Table) NewIndex(name string, columns ...*Column) Index {
	index := Index{
		DBName:    table.DBName,
		TableName: table.Name,
		Table:     table,
		Name:      name,
		Columns:   columns,
	}
	return index
}

// PushIndex push an index instance to the table indexes
func (table *Table) PushIndex(index *Index) *Table {
	table.IndexMap[index.Name] = index
	table.Indexes = append(table.Indexes, index)
	return table
}

// NewColumn create a new column intstance
func (table *Table) NewColumn(name string) Column {
	column := Column{
		DBName:            table.DBName,
		TableName:         table.Name,
		Table:             table,
		Name:              name,
		Length:            nil,
		OctetLength:       nil,
		Precision:         nil,
		Scale:             nil,
		DatetimePrecision: nil,
		Charset:           nil,
		Collation:         nil,
		Key:               nil,
		Extra:             nil,
		Comment:           nil,
	}
	return column
}

// PushColumn push a column instance to the table columns
func (table *Table) PushColumn(column *Column) *Table {
	table.ColumnMap[column.Name] = column
	table.Columns = append(table.Columns, column)
	return table
}

// HasColumn checking if the given name column exists
func (table *Table) HasColumn(name string) bool {
	_, has := table.ColumnMap[name]
	return has
}

// HasIndex checking if the given name index exists
func (table *Table) HasIndex(name string) bool {
	_, has := table.IndexMap[name]
	return has
}

// AddCommand Add a new command to the table.
//
// The commands must be:
//    AddColumn(column *Column)    for adding a column
//    ModifyColumn(column *Column) for modifying a colu
//    RenameColumn(old string,new string)  for renaming a column
//    DropColumn(name string)  for dropping a column
//    CreateIndex(index *Index) for creating a index
//    DropIndex( name string) for  dropping a index
//    RenameIndex(old string,new string)  for renaming a index
func (table *Table) AddCommand(name string, params ...interface{}) {
	table.Commands = append(table.Commands, &Command{
		Name:   name,
		Params: params,
	})
}
