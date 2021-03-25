package query

// Insert Insert new records into the database.
func (builder *Builder) Insert() {
}

// MustInsert Insert new records into the database.
func (builder *Builder) MustInsert() {
}

// InsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) InsertOrIgnore() {
}

// MustInsertOrIgnore Insert new records into the database while ignoring errors.
func (builder *Builder) MustInsertOrIgnore() {
}

// InsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) InsertGetID() {
}

// MustInsertGetID Insert a new record and get the value of the primary key.
func (builder *Builder) MustInsertGetID() {
}

// InsertUsing Insert new records into the table using a subquery.
func (builder *Builder) InsertUsing() {
}

// MustInsertUsing Insert new records into the table using a subquery.
func (builder *Builder) MustInsertUsing() {
}
