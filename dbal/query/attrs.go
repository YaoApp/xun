package query

// From set the table which the query is targeting.
func (builder *Builder) From(fullname string) Query {
	builder.SetAttrFrom(fullname)
	return builder
}

// FromSub Makes "from" fetch from a subquery.
func (builder *Builder) FromSub() {
}

// FromRaw Add a raw from clause to the query.
func (builder *Builder) FromRaw() {
}

// SetAttrFrom set the From attribute
func (builder *Builder) SetAttrFrom(fullname string) {
	builder.Attr.From = builder.Name(fullname)
}

// AddColumn add the query column to the builder
func (builder *Builder) AddColumn(fullname string) {
	builder.Attr.Columns = append(builder.Attr.Columns, builder.Name(fullname))
}
