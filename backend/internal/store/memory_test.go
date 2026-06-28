package store

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_EventTypes(t *testing.T) {
	s := NewMemoryStore()

	et, err := models.NewEventType("et-1", "Meeting", "Description", 30)
	require.NoError(t, err)

	s.CreateEventType(et)

	list := s.ListEventTypes()
	require.Len(t, list, 1)
	assert.Equal(t, et, list[0])

	got, ok := s.GetEventType("et-1")
	require.True(t, ok)
	assert.Equal(t, et, got)

	updated, err := models.NewEventType("et-1", "Updated", "New description", 60)
	require.NoError(t, err)
	assert.True(t, s.UpdateEventType(updated))

	got, ok = s.GetEventType("et-1")
	require.True(t, ok)
	assert.Equal(t, updated, got)

	assert.True(t, s.DeleteEventType("et-1"))
	assert.False(t, s.DeleteEventType("et-1"))
	_, ok = s.GetEventType("et-1")
	assert.False(t, ok)
}

func TestMemoryStore_Bookings(t *testing.T) {
	s := NewMemoryStore()
	now := time.Date(2026, 6, 28, 10, 0, 0, 0, time.UTC)
	start := now.Add(time.Hour)

	b1, err := models.NewBooking("b-1", "et-1", "Alice", start, 30*time.Minute, now)
	require.NoError(t, err)

	require.NoError(t, s.CreateBooking(b1))

	list := s.ListBookings()
	require.Len(t, list, 1)
	assert.Equal(t, b1, list[0])

	got, ok := s.GetBooking("b-1")
	require.True(t, ok)
	assert.Equal(t, b1, got)

	b2, err := models.NewBooking("b-2", "et-1", "Bob", start, 30*time.Minute, now)
	require.NoError(t, err)
	assert.ErrorIs(t, s.CreateBooking(b2), ErrSlotTaken)

	otherStart := start.Add(30 * time.Minute)
	b3, err := models.NewBooking("b-3", "et-1", "Charlie", otherStart, 30*time.Minute, now)
	require.NoError(t, err)
	require.NoError(t, s.CreateBooking(b3))

	assert.True(t, s.DeleteBooking("b-1"))
	assert.False(t, s.DeleteBooking("b-1"))
}

func TestMemoryStore_CreateBooking_Race(t *testing.T) {
	s := NewMemoryStore()
	now := time.Date(2026, 6, 28, 10, 0, 0, 0, time.UTC)
	start := now.Add(time.Hour)

	var wg sync.WaitGroup
	success := make(chan int, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			b, err := models.NewBooking(fmt.Sprintf("b-%d", i), "et-1", "Guest", start, 30*time.Minute, now)
			require.NoError(t, err)
			if err := s.CreateBooking(b); err == nil {
				success <- 1
			}
		}(i)
	}

	wg.Wait()
	close(success)

	count := 0
	for range success {
		count++
	}
	assert.Equal(t, 1, count)
}
