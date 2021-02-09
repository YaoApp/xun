package grammar

// NewTable create a grammar table
func NewTable(name string, dbname string) Table {
	return Table{
		DBName:    dbname,
		Name:      name,
		Columns:   []*Column{},
		ColumnMap: map[string]*Column{},
		Indexes:   []*Index{},
		IndexMap:  map[string]*Index{},
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

// AddIndex add a index to the table
func (table *Table) AddIndex(index *Index) *Table {
	table.IndexMap[index.Name] = index
	table.Indexes = append(table.Indexes, index)
	return table
}

// NewColumn create a new column intstance
func (table *Table) NewColumn(name string) Column {
	column := Column{
		DBName:    table.DBName,
		TableName: table.Name,
		Table:     table,
		Name:      name,
	}
	return column
}

// AddColumn add a column to the table
func (table *Table) AddColumn(column *Column) *Table {
	table.ColumnMap[column.Name] = column
	table.Columns = append(table.Columns, column)
	return table
}
