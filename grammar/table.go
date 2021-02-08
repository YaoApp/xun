package grammar

// SetDefaultEngine set the default engine
func (table *Table) SetDefaultEngine(engine string) {
	if table.Engine == "" {
		table.Engine = engine
	}
}

// SetDefaultCollation set the default collation
func (table *Table) SetDefaultCollation(collation string) {
	if table.Collation == "" {
		table.Collation = collation
	}
}

// SetDefaultCharset set the default charset
func (table *Table) SetDefaultCharset(charset string) {
	if table.Charset == "" {
		table.Charset = charset
	}
}
