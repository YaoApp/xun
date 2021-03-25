package query

// Where Add a basic where clause to the query.
func (builder *Builder) Where() {
}

// OrWhere Add an "or where" clause to the query.
func (builder *Builder) OrWhere() {
}

// WhereJSONContains Add a "where JSON contains" clause to the query.
func (builder *Builder) WhereJSONContains() {
}

// OrWhereJSONContains Add an "or where JSON contains" clause to the query.
func (builder *Builder) OrWhereJSONContains() {
}

// WhereJSONDoesntContain Add a "where JSON not contains" clause to the query.
func (builder *Builder) WhereJSONDoesntContain() {
}

// OrWhereJSONDoesntContain Add an "or where JSON not contains" clause to the query.
func (builder *Builder) OrWhereJSONDoesntContain() {
}

// WhereJSONLength Add a "where JSON length" clause to the query.
func (builder *Builder) WhereJSONLength() {
}

// OrWhereJSONLength Add an "or where JSON length" clause to the query.
func (builder *Builder) OrWhereJSONLength() {
}

// WhereBetween Add a where between statement to the query.
func (builder *Builder) WhereBetween() {
}

// OrWhereBetween Add an or where between statement to the query.
func (builder *Builder) OrWhereBetween() {
}

// WhereNotBetween Add a where not between statement to the query.
func (builder *Builder) WhereNotBetween() {
}

// OrWhereNotBetween Add an or where not between statement using columns to the query.
func (builder *Builder) OrWhereNotBetween() {
}

// WhereIn Add a "where in" clause to the query.
func (builder *Builder) WhereIn() {
}

// OrWhereIn Add an "or where in" clause to the query.
func (builder *Builder) OrWhereIn() {
}

// WhereNotIn Add a "where not in" clause to the query.
func (builder *Builder) WhereNotIn() {
}

// OrWhereNotIn Add an "or where not in" clause to the query.
func (builder *Builder) OrWhereNotIn() {
}

// WhereNull Add a "where null" clause to the query.
func (builder *Builder) WhereNull() {
}

// OrWhereNull Add an "or where null" clause to the query.
func (builder *Builder) OrWhereNull() {
}

// WhereNull Add a "where not null" clause to the query.
func (builder *Builder) whereNotNull() {
}

// OrWhereNotNull Add an "or where not null" clause to the query.
func (builder *Builder) OrWhereNotNull() {
}

// WhereDate Add a "where date" statement to the query.
func (builder *Builder) WhereDate() {
}

// OrWhereDate Add an "or where date" statement to the query.
func (builder *Builder) OrWhereDate() {
}

// WhereYear Add a "where year" statement to the query.
func (builder *Builder) WhereYear() {
}

// OrWhereYear Add an "or where year" statement to the query.
func (builder *Builder) OrWhereYear() {
}

// WhereMonth Add a "where month" statement to the query.
func (builder *Builder) WhereMonth() {
}

// OrWhereMonth Add an "or where month" statement to the query.
func (builder *Builder) OrWhereMonth() {
}

// WhereDay Add a "where day" statement to the query.
func (builder *Builder) WhereDay() {
}

// OrWhereDay Add an "or where day" statement to the query.
func (builder *Builder) OrWhereDay() {
}

// WhereTime Add a "where time" statement to the query.
func (builder *Builder) WhereTime() {
}

// OrWhereTime Add an "or where time" statement to the query.
func (builder *Builder) OrWhereTime() {
}

// WhereColumn Add a "where" clause comparing two columns to the query.
func (builder *Builder) WhereColumn() {
}

// OrWhereColumn Add an "or where" clause comparing two columns to the query.
func (builder *Builder) OrWhereColumn() {
}

// WhereExists Add an exists clause to the query.
func (builder *Builder) WhereExists() {
}

// OrWhereExists Add an or exists clause to the query.
func (builder *Builder) OrWhereExists() {
}

// WhereNotExists  Add a where not exists clause to the query.
func (builder *Builder) WhereNotExists() {
}

// OrWhereNotExists Add a where not exists clause to the query.
func (builder *Builder) OrWhereNotExists() {
}

// WhereRaw Add a basic where clause to the query.
func (builder *Builder) WhereRaw() {
}

// OrWhereRaw Add an "or where" clause to the query.
func (builder *Builder) OrWhereRaw() {
}
