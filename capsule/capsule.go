package capsule

// New Create a database manager instance.
func New() *Manager {
	return &Manager{
		pool: []*Connection{},
	}
}

// AddConnection Register a connection with the manager.
func (manager *Manager) AddConnection() {
}

// GetConnection Get a registered connection instance.
func (manager *Manager) GetConnection() {
}

// Schema Get a schema builder instance.
func (manager *Manager) Schema() {
}

// Query Get a fluent query builder instance.
func (manager *Manager) Query() {
}
