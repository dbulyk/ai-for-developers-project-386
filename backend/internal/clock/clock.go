package clock

import "time"

// Clock provides the current time. It is implemented by RealClock and
// MockClock so production code and tests can control time deterministically.
type Clock interface {
	Now() time.Time
}

// RealClock returns the actual current time.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time {
	return time.Now()
}
