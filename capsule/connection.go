package capsule

import (
	"context"
	"time"
)

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (conn *Connection) Ping(timeout time.Duration) (err error) {

	done := make(chan bool, 1)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		err = conn.DB.PingContext(ctx)
		done <- true
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		break
	case <-done:
		break
	}

	return err
}
