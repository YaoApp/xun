package query

// Query The database Query interface
type Query interface {
	Where()
	Join()
}

// Builder the dbal query builder
type Builder struct{ Query }
