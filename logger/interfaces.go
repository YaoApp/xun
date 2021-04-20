package logger

import "time"

// Writer logger writer
type Writer interface {
	Trace(trace ...string) Writer
	TimeCost(start time.Time)
	Write()
}

// Empty the empty writer
type Empty struct{}

// TimeCost empty TimeCost
func (empty Empty) TimeCost(start time.Time) {
}

// TimeCost empty Write
func (empty Empty) Write() {
}

// Trace empty trace.
func (empty Empty) Trace(trace ...string) Writer {
	return empty
}
