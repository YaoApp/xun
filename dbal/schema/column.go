package schema

// Drop mark as dropped for the index
func (column *Column) Drop() {
	column.Dropped = true
}

// Rename mark as renamed with the given name for the index
func (column *Column) Rename(new string) {
	column.Newname = new
}

// Unique set as index
func (column *Column) Unique() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "unique"
	column.Table.AddIndex(index)
	return column
}

// Primary set as primary key
func (column *Column) Primary() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "primary"
	column.Table.AddIndex(index)
	column.Column.Primary = true
	return column
}

// Index set as index key
func (column *Column) Index() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "index"
	column.Table.AddIndex(index)
	column.Column.Primary = true
	return column
}
