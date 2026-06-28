package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/slots"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPublicEventTypesHandler(now time.Time) (*PublicEventTypesHandler, *store.MemoryStore, *clock.MockClock) {
	s := store.NewMemoryStore().(*store.MemoryStore)
	c := clock.NewMockClock(now)
	tz := time.UTC
	h := NewPublicEventTypesHandler(s, c, tz)
	return h, s, c
}

func TestPublicEventTypes_List(t *testing.T) {
	h, s, _ := setupPublicEventTypesHandler(time.Now())
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	et := models.EventType{ID: "et-1", Name: "Consultation", Description: "Quick chat", DurationMinutes: 30}
	s.CreateEventType(et)

	req := httptest.NewRequest(http.MethodGet, "/public/event-types", http.NoBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got []models.EventType
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	assert.Equal(t, et, got[0])
}

func TestPublicEventTypes_GetSlots(t *testing.T) {
	// 2026-06-29 07:00 UTC is Monday 10:00 Europe/Moscow -> first slot 09:00 MSK = 06:00 UTC.
	moscow, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	now := time.Date(2026, 6, 29, 7, 0, 0, 0, time.UTC)
	h, s, _ := setupPublicEventTypesHandler(now)
	h.tz = moscow

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	et := models.EventType{ID: "et-1", Name: "Consultation", DurationMinutes: 60}
	s.CreateEventType(et)

	req := httptest.NewRequest(http.MethodGet, "/public/event-types/et-1/slots", http.NoBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got []slots.AvailableDay
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 14)

	// Today is Monday: 09:00-18:00 MSK. 07:00 UTC -> 10:00 MSK, so 09:00 MSK slot has passed.
	// Expected slots: 10:00..17:00 MSK -> 07:00..14:00 UTC (8 slots).
	require.Len(t, got[0].Slots, 7)
	assert.Equal(t, "2026-06-29T08:00:00Z", got[0].Slots[0])
	assert.Equal(t, "2026-06-29T14:00:00Z", got[0].Slots[len(got[0].Slots)-1])
}

func TestPublicEventTypes_GetSlots_NotFound(t *testing.T) {
	h, _, _ := setupPublicEventTypesHandler(time.Now())
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/public/event-types/missing/slots", http.NoBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "not_found")
}
