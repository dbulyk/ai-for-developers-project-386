package clock

import (
	"sync"
	"time"
)

// MockClock is a test-only Clock with a mutable, concurrency-safe time.
type MockClock struct {
	mu sync.Mutex
	t  time.Time
}

// NewMockClock creates a MockClock set to the provided time.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{t: t}
}

// Now returns the current mock time.
func (m *MockClock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.t
}

// Advance moves the mock time forward by d.
func (m *MockClock) Advance(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.t = m.t.Add(d)
}

// Set sets the mock time to t.
func (m *MockClock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.t = t
}
