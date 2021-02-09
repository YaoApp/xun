package schema

var columnTypes = map[string]string{
	"string": "VARCHAR",
}

var fieldTypes = map[string]string{
	"VARCHAR": "string",
}

// Unique set as index
func (column *Column) Unique() *Column {
	index := column.Table.NewIndex(column.Name, column)
	index.Type = "unique"
	column.Table.AddIndex(index)
	return column
}

// BigInteger Create a new auto-incrementing big integer (8-byte) column on the table.
func (table *Table) BigInteger() {}

// String Create a new string column on the table.
func (table *Table) String(name string, length int) *Column {
	column := table.NewColumn(name)
	column.Length = length
	column.Type = "string"
	table.AddColumn(column)
	return column
}

// Drop mark as dropped for the index
func (column *Column) Drop() {
	column.Dropped = true
}

// Rename mark as renamed with the given name for the index
func (column *Column) Rename(new string) {
	column.Newname = new
}
