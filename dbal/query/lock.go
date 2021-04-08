package query

// SharedLock Share lock the selected rows in the table.
func (builder *Builder) SharedLock() Query {
	return builder.Lock("share")
}

// LockForUpdate Lock the selected rows in the table for updating.
func (builder *Builder) LockForUpdate() Query {
	return builder.Lock("update")
}

// Lock Lock the selected rows in the table.
func (builder *Builder) Lock(value interface{}) Query {
	builder.Query.Lock = value
	if builder.Query.Lock != nil {
		builder.UseWrite()
	}
	return builder
}
