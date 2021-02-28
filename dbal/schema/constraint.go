package schema

// AddUniqueConstraint add a unique coustraint
func (table *Table) AddUniqueConstraint(name string, columnNames ...string) {
}

// GetUniqueConstraint Returns the unique constraint with the given name.
func (table *Table) GetUniqueConstraint(name string) {
}

// DropUniqueConstraint  Removes the unique constraint with the given name.
func (table *Table) DropUniqueConstraint(name string) {
}
