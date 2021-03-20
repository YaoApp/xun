package schema

import "github.com/yaoapp/xun/dbal"

// addColumnCommand add a new command that adding a column
func (table *Table) addColumnCommand(column *dbal.Column, success func(), fail func()) {
	table.AddCommand("AddColumn", success, fail, column)
}

// changeColumnCommand add a new command that modifing a column
func (table *Table) changeColumnCommand(column *dbal.Column, success func(), fail func()) {
	table.AddCommand("ChangeColumn", success, fail, column)
}

// renameColumnCommand add a new command that renaming a column
func (table *Table) renameColumnCommand(old string, new string, success func(), fail func()) {
	table.AddCommand("RenameColumn", success, fail, old, new)
}

// dropColumnCommand add a new command that dropping a column
func (table *Table) dropColumnCommand(name string, success func(), fail func()) {
	table.AddCommand("DropColumn", success, fail, name)
}

// createIndexCommand add a new command that creating a index
func (table *Table) createIndexCommand(index *dbal.Index, success func(), fail func()) {
	table.AddCommand("CreateIndex", success, fail, index)
}

// createPrimaryCommand add a new command that creating the primary key
func (table *Table) createPrimaryCommand(primary *dbal.Primary, success func(), fail func()) {
	table.AddCommand("CreatePrimary", success, fail, primary)
}

// dropPrimaryCommand add a new command drop the primary key
func (table *Table) dropPrimaryCommand(primary *Primary, success func(), fail func()) {
	if primary == nil {
		success()
	}
	table.AddCommand("DropPrimary", success, fail, primary.Name, primary.Columns)
}

// dropIndexCommand add a new command that dropping a index
func (table *Table) dropIndexCommand(name string, success func(), fail func()) {
	table.AddCommand("DropIndex", success, fail, name)
}

// renameIndexCommand add a new command that renaming a index
func (table *Table) renameIndexCommand(old string, new string, success func(), fail func()) {
	table.AddCommand("RenameIndex", success, fail, old, new)
}
