package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPublicBookingsHandler(now time.Time) (*PublicBookingsHandler, *store.MemoryStore, *clock.MockClock) {
	s := store.NewMemoryStore().(*store.MemoryStore)
	c := clock.NewMockClock(now)
	tz := time.UTC
	h := NewPublicBookingsHandler(s, c, tz)
	return h, s, c
}

func TestPublicBookings_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
		h, s, _ := setupPublicBookingsHandler(now)
		h.tz = mustLocation(t, "Europe/Moscow")

		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})

		body := `{"eventTypeId":"et-1","guestName":"Alice","startTime":"2026-06-29T07:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var got models.Booking
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "et-1", got.EventTypeID)
		assert.Equal(t, "Alice", got.GuestName)
		assert.Equal(t, time.Date(2026, 6, 29, 7, 0, 0, 0, time.UTC), got.StartTime)
		assert.Equal(t, time.Date(2026, 6, 29, 8, 0, 0, 0, time.UTC), got.EndTime)
		assert.False(t, got.CreatedAt.IsZero())

		stored, ok := s.GetBooking(got.ID)
		require.True(t, ok)
		assert.Equal(t, got, stored)
	})

	t.Run("event type not found", func(t *testing.T) {
		h, _, _ := setupPublicBookingsHandler(time.Now())
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		body := `{"eventTypeId":"missing","guestName":"Alice","startTime":"2026-06-29T07:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "not_found")
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _, _ := setupPublicBookingsHandler(time.Now())
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid_json")
	})

	t.Run("invalid start time format", func(t *testing.T) {
		h, s, _ := setupPublicBookingsHandler(time.Now())
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})

		body := `{"eventTypeId":"et-1","guestName":"Alice","startTime":"not-a-time"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "start_time_invalid")
	})

	t.Run("slot not in available window", func(t *testing.T) {
		now := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
		h, s, _ := setupPublicBookingsHandler(now)
		h.tz = mustLocation(t, "Europe/Moscow")

		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})

		// 2026-06-29T05:00:00Z is 08:00 MSK (before work hours).
		body := `{"eventTypeId":"et-1","guestName":"Alice","startTime":"2026-06-29T05:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "slot_invalid")
	})

	t.Run("slot already taken returns 409", func(t *testing.T) {
		now := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
		h, s, _ := setupPublicBookingsHandler(now)
		h.tz = mustLocation(t, "Europe/Moscow")

		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})
		existing := models.Booking{
			ID:          "b-1",
			EventTypeID: "et-1",
			GuestName:   "Bob",
			StartTime:   time.Date(2026, 6, 29, 11, 0, 0, 0, time.UTC),
			EndTime:     time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC),
			CreatedAt:   now,
		}
		require.NoError(t, s.CreateBooking(existing))

		body := `{"eventTypeId":"et-1","guestName":"Alice","startTime":"2026-06-29T11:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "slot_taken")
	})

	t.Run("guest name required", func(t *testing.T) {
		now := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
		h, s, _ := setupPublicBookingsHandler(now)
		h.tz = mustLocation(t, "Europe/Moscow")

		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})

		body := `{"eventTypeId":"et-1","guestName":"","startTime":"2026-06-29T07:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "guest_name_required")
	})
}

func TestPublicBookings_Create_Race(t *testing.T) {
	now := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
	h, s, _ := setupPublicBookingsHandler(now)
	h.tz = mustLocation(t, "Europe/Moscow")

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	s.CreateEventType(models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60})

	body := `{"eventTypeId":"et-1","guestName":"Alice","startTime":"2026-06-29T07:00:00Z"}`

	var wg sync.WaitGroup
	var created int
	var mu sync.Mutex

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/public/bookings", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code == http.StatusCreated {
				mu.Lock()
				created++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, 1, created)
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	require.NoError(t, err)
	return loc
}
