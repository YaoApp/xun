package query

// Select Set the columns to be selected.
func (builder *Builder) Select() {
}

// SelectSub Add a subselect expression to the query.
func (builder *Builder) SelectSub() {
}

// SelectRaw Add a new "raw" select expression to the query.
func (builder Builder) SelectRaw() {
}

// Distinct Force the query to only return distinct results.
func (builder Builder) Distinct() {
}
