package service

import "time"

// Timer type
type Timer struct {
	StartTime time.Time
}

// NewTimer returns a new Timer instance
func NewTimer() *Timer {
	return &Timer{StartTime: time.Now()}
}

// to return Time in seconds
func (t *Timer) durationSeconds() float64 {
	return float64(time.Since(t.StartTime).Seconds())
}
