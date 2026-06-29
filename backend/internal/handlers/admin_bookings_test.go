package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminBookingsHandler(now time.Time) (*AdminBookingsHandler, *store.MemoryStore, *clock.MockClock) {
	s := store.NewMemoryStore().(*store.MemoryStore)
	c := clock.NewMockClock(now)
	h := NewAdminBookingsHandler(s, c)
	return h, s, c
}

func TestAdminBookings_List(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)

	t.Run("empty list", func(t *testing.T) {
		h, _, _ := setupAdminBookingsHandler(now)
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodGet, "/admin/bookings", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, "[]", rec.Body.String())
	})

	t.Run("filters past bookings and sorts upcoming", func(t *testing.T) {
		h, s, _ := setupAdminBookingsHandler(now)
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		past := models.Booking{
			ID:          "b-past",
			EventTypeID: "et-1",
			GuestName:   "Past",
			StartTime:   now.Add(-1 * time.Hour),
			EndTime:     now.Add(-30 * time.Minute),
			CreatedAt:   now,
		}
		future1 := models.Booking{
			ID:          "b-future1",
			EventTypeID: "et-1",
			GuestName:   "Future1",
			StartTime:   now.Add(2 * time.Hour),
			EndTime:     now.Add(3 * time.Hour),
			CreatedAt:   now,
		}
		future2 := models.Booking{
			ID:          "b-future2",
			EventTypeID: "et-1",
			GuestName:   "Future2",
			StartTime:   now.Add(1 * time.Hour),
			EndTime:     now.Add(2 * time.Hour),
			CreatedAt:   now,
		}

		require.NoError(t, s.CreateBooking(past))
		require.NoError(t, s.CreateBooking(future1))
		require.NoError(t, s.CreateBooking(future2))

		req := httptest.NewRequest(http.MethodGet, "/admin/bookings", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got []models.Booking
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		require.Len(t, got, 2)
		assert.Equal(t, "b-future2", got[0].ID)
		assert.Equal(t, "b-future1", got[1].ID)
	})
}

func TestAdminBookings_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, s, _ := setupAdminBookingsHandler(time.Now())
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		b := models.Booking{
			ID:          "b-1",
			EventTypeID: "et-1",
			GuestName:   "Alice",
			StartTime:   time.Now().Add(time.Hour),
			EndTime:     time.Now().Add(2 * time.Hour),
			CreatedAt:   time.Now(),
		}
		require.NoError(t, s.CreateBooking(b))

		req := httptest.NewRequest(http.MethodDelete, "/admin/bookings/b-1", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())

		_, ok := s.GetBooking("b-1")
		assert.False(t, ok)
	})

	t.Run("not found", func(t *testing.T) {
		h, _, _ := setupAdminBookingsHandler(time.Now())
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodDelete, "/admin/bookings/missing", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "not_found")
	})
}
