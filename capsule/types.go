package capsule

// Manager The database manager
type Manager struct {
	pool []*Connection
}

// Connection The database connection
type Connection struct{}
