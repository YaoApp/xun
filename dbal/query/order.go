package query

// OrderBy Add an "order by" clause to the query.
func (builder *Builder) OrderBy() {
}

// Latest Add an "order by" clause for a timestamp to the query.
func (builder *Builder) Latest() {
}

// Oldest Add an "order by" clause for a timestamp to the query.
func (builder *Builder) Oldest() {
}

// InRandomOrder Put the query's results in random order.
func (builder *Builder) InRandomOrder() {
}

// Reorder Remove all existing orders and optionally add a new order.
func (builder *Builder) Reorder() {
}

// OrderByRaw Add a raw "order by" clause to the query.
func (builder *Builder) OrderByRaw() {
}
