package schema

import (
	"github.com/yaoapp/xun/dbal"
)

// GetPrimary get the table primary key instance
func (table *Table) GetPrimary() *Primary {
	return table.Primary
}

// AddPrimary Indicate that the given column should be a primary index.
func (table *Table) AddPrimary(columnNames ...string) {
	table.addPrimaryWithName("PRIMARY", columnNames...)
}

// AddPrimaryWithName Indicate that the given column should be a primary index.
func (table *Table) addPrimaryWithName(name string, columnNames ...string) {
	columns := []*dbal.Column{}
	for _, columnName := range columnNames {
		column := table.GetColumn(columnName)
		column.NotNull()
		column.Column.Primary = true
		columns = append(columns, column.Column)
	}
	primary := &Primary{
		Primary: table.NewPrimary(name, columns...),
		Table:   table,
	}

	table.Primary = primary
	table.CreatePrimaryCommand(primary.Primary, nil, func() {
		table.Primary = nil
	})
}

// DropPrimary Indicate that dropping the primary index
func (table *Table) DropPrimary() {
	primary := table.GetPrimary()
	table.DropPrimaryCommand(primary, func() {
		table.Primary = nil
	}, nil)
}
