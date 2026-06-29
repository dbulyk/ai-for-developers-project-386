package handlers

import (
	"net/http"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/slots"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// PublicEventTypesHandler exposes public read operations for event types and
// their available slots.
type PublicEventTypesHandler struct {
	store store.Store
	clock clock.Clock
	tz    *time.Location
}

// NewPublicEventTypesHandler creates a new PublicEventTypesHandler.
func NewPublicEventTypesHandler(s store.Store, c clock.Clock, tz *time.Location) *PublicEventTypesHandler {
	return &PublicEventTypesHandler{
		store: s,
		clock: c,
		tz:    tz,
	}
}

// RegisterRoutes mounts the public event types routes on the given router.
func (h *PublicEventTypesHandler) RegisterRoutes(r chi.Router) {
	r.Get("/public/event-types", h.List)
	r.Get("/public/event-types/{id}/slots", h.GetSlots)
}

// List returns all event types.
func (h *PublicEventTypesHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.ListEventTypes())
}

// GetSlots returns available slots for the requested event type.
func (h *PublicEventTypesHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id_required", "id is required")
		return
	}

	et, ok := h.store.GetEventType(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "event type not found")
		return
	}

	bookings := h.store.ListBookings()
	taken := make([]time.Time, 0, len(bookings))
	for _, b := range bookings {
		if b.EventTypeID == id {
			taken = append(taken, b.StartTime)
		}
	}

	duration := time.Duration(et.DurationMinutes) * time.Minute
	available := slots.Generate(id, duration, h.tz, h.clock, taken)
	writeJSON(w, http.StatusOK, available)
}
