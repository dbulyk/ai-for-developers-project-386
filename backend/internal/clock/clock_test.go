package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/stretchr/testify/assert"
)

func TestRealClock_Now(t *testing.T) {
	rc := clock.RealClock{}

	before := time.Now()
	got := rc.Now()
	after := time.Now()

	assert.WithinDuration(t, before, got, time.Second)
	assert.WithinDuration(t, after, got, time.Second)
}

func TestMockClock_Now(t *testing.T) {
	start := time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC)
	mc := clock.NewMockClock(start)

	assert.Equal(t, start, mc.Now())
}

func TestMockClock_Advance(t *testing.T) {
	start := time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC)
	mc := clock.NewMockClock(start)

	mc.Advance(30 * time.Minute)

	assert.Equal(t, start.Add(30*time.Minute), mc.Now())
}

func TestMockClock_Set(t *testing.T) {
	start := time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC)
	mc := clock.NewMockClock(start)
	newTime := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)

	mc.Set(newTime)

	assert.Equal(t, newTime, mc.Now())
}

func TestMockClock_Concurrency(t *testing.T) {
	start := time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC)
	mc := clock.NewMockClock(start)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = mc.Now()
			mc.Advance(time.Minute)
			mc.Set(start)
		}()
	}
	wg.Wait()

	// The primary goal of this test is to ensure there are no data races.
	// The exact final value is non-deterministic because Advance/Set sequences
	// from different goroutines can interleave.
	_ = mc.Now()
}
