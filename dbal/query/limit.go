package query

// Skip Alias to set the "offset" value of the query.
func (builder *Builder) Skip(value int) Query {
	return builder.Offset(value)
}

// Offset Set the "offset" value of the query.
func (builder *Builder) Offset(value int) Query {
	if value < 0 {
		value = 0
	}
	if len(builder.Query.Unions) > 0 {
		builder.Query.UnionOffset = value
	} else {
		builder.Query.Offset = value
	}
	return builder
}

// Take Alias to set the "limit" value of the query.
func (builder *Builder) Take(value int) Query {
	return builder.Limit(value)
}

// Limit Set the "limit" value of the query.
func (builder *Builder) Limit(value int) Query {
	if value < 0 {
		return builder
	}
	if len(builder.Query.Unions) > 0 {
		builder.Query.UnionLimit = value
	} else {
		builder.Query.Limit = value
	}
	return builder
}
