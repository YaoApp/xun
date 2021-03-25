package query

// Update Update records in the database.
func (builder *Builder) Update() {
}

// MustUpdate Update records in the database.
func (builder *Builder) MustUpdate() {
}

// UpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) UpdateOrInsert() {
}

// MustUpdateOrInsert Insert or update a record matching the attributes, and fill it with values.
func (builder *Builder) MustUpdateOrInsert() {
}

// Upsert new records or update the existing ones.
func (builder *Builder) Upsert() {
}

// MustUpsert new records or update the existing ones.
func (builder *Builder) MustUpsert() {
}

// Increment Increment a column's value by a given amount.
func (builder *Builder) Increment() {
}

// MustIncrement Increment a column's value by a given amount.
func (builder *Builder) MustIncrement() {
}

// Decrement Decrement a column's value by a given amount.
func (builder *Builder) Decrement() {
}

// MustDecrement Decrement a column's value by a given amount.
func (builder *Builder) MustDecrement() {
}
