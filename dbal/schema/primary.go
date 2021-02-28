package schema

import "github.com/yaoapp/xun/dbal"

// GetPrimary get the table primary key instance
func (table *Table) GetPrimary() *Primary {
	if table.Primary == nil && table.Table.Primary != nil {
		return &Primary{
			Primary: *table.Table.Primary,
		}
	}
	return table.Primary
}

// CreatePrimary Indicate that the given column should be a primary index.
func (table *Table) CreatePrimary(columnNames ...string) {
	table.CreatePrimaryWithName("PRIMARY", columnNames...)
}

// CreatePrimaryWithName Indicate that the given column should be a primary index.
func (table *Table) CreatePrimaryWithName(name string, columnNames ...string) {
	columns := []*dbal.Column{}
	for _, columnName := range columnNames {
		column := table.GetColumn(columnName)
		column.NotNull()
		column.Column.Primary = true
		columns = append(columns, &column.Column)
	}
	primary := Primary{
		Primary: table.NewPrimary(name, columns...),
		Table:   table,
	}
	table.CreatePrimaryCommand(&primary.Primary)
	table.onChange("CreatePrimary", primary)
}

// DropPrimary Indicate that dropping the primary index
func (table *Table) DropPrimary() {
	table.DropPrimaryCommand()
	table.onChange("DropPrimary", table.Primary)
	table.Primary = nil
}
