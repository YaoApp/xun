package query

// From set the table which the query is targeting.
func (builder *Builder) From(name string) Query {
	builder.SetAttrFrom(name)
	return builder
}

// FromSub Makes "from" fetch from a subquery.
func (builder *Builder) FromSub() {
}

// FromRaw Add a raw from clause to the query.
func (builder *Builder) FromRaw() {
}

// SetAttrFrom set the From attribute
func (builder *Builder) SetAttrFrom(name string) {
	builder.Attr.From = builder.NewTable(name)
}
