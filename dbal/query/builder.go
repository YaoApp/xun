package query

// New Get a fluent query builder instance.
func New(conn *Connection) Query {
	builder := NewBuilder(conn)
	return &builder
}

// NewBuilder create new query builder instance
func NewBuilder(conn *Connection) Builder {
	return Builder{
		Conn: conn,
	}
}

// Where Add a basic where clause to the query.
func (builder *Builder) Where() {
	// fmt.Printf("\nWhere DBAL: \n===\n%#v\n===\n", builder.Conn)
}

// Join Add a join clause to the query.
func (builder *Builder) Join() {}
