package schema

// CreateIndex Indicate that the given index should be created.
func (table *Table) CreateIndex(key string, columnNames ...string) {
	columns := []*Column{}
	for _, name := range columnNames {
		columns = append(columns, table.GetColumn(name))
	}
	index := table.NewIndex(key, columns...)
	index.Type = "index"
	table.CreateIndexCommand(&index.Index)
	table.onChange("CreateIndex", index)
}

// CreateUnique Indicate that the given unique index should be created.
func (table *Table) CreateUnique(key string, columnNames ...string) {
	columns := []*Column{}
	for _, name := range columnNames {
		columns = append(columns, table.GetColumn(name))
	}
	index := table.NewIndex(key, columns...)
	index.Type = "unique"
	table.CreateIndexCommand(&index.Index)
	table.onChange("CreateIndex", index)
}

// DropIndex Indicate that the given indexes should be dropped.
func (table *Table) DropIndex(key ...string) {
	for _, n := range key {
		table.DropIndexCommand(n)
	}
	table.onChange("DropIndex", key)
}

// RenameIndex Indicate that the given indexes should be renamed.
func (table *Table) RenameIndex(old string, new string) *Index {
	table.RenameIndexCommand(old, new)
	index := table.GetIndex(old)
	index.Name = new
	table.onChange("RenameIndex", old, new)
	return index
}
