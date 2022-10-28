package capsule

import (
	"fmt"
	"math/rand"
	"time"
)

// RandPrimary rand select primary connection
func (pool *Pool) RandPrimary() (*Connection, error) {

	length := len(pool.Primary)
	if length == 0 {
		return nil, fmt.Errorf("the primary connection was empty")
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	i := r.Intn(length)
	return pool.Primary[i], nil
}

// RandReadOnly rand select primary connection
func (pool *Pool) RandReadOnly() (*Connection, error) {
	length := len(pool.Readonly)
	if length == 0 {
		return pool.RandPrimary()
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	i := r.Intn(length)
	return pool.Readonly[i], nil
}
