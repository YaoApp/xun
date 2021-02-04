package query

import "fmt"

// Table Get a fluent query builder instance.
func Table() Query {
	builder := NewBuilder()
	return &builder
}

// NewBuilder create new query builder instance
func NewBuilder() Builder {
	return Builder{}
}

// Where Add a basic where clause to the query.
func (builder *Builder) Where() { fmt.Printf("DBAL WHERE\n") }

// Join Add a join clause to the query.
func (builder *Builder) Join() { fmt.Printf("DBAL JOIN\n") }
