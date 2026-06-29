package handlers

import (
	"net/http"
	"sort"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// AdminBookingsHandler exposes admin operations for bookings.
type AdminBookingsHandler struct {
	store store.Store
	clock clock.Clock
}

// NewAdminBookingsHandler creates a new AdminBookingsHandler.
func NewAdminBookingsHandler(s store.Store, c clock.Clock) *AdminBookingsHandler {
	return &AdminBookingsHandler{
		store: s,
		clock: c,
	}
}

// RegisterRoutes mounts the admin bookings routes on the given router.
func (h *AdminBookingsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/admin/bookings", h.List)
	r.Delete("/admin/bookings/{id}", h.Delete)
}

// List returns all upcoming bookings sorted by start time ascending.
func (h *AdminBookingsHandler) List(w http.ResponseWriter, r *http.Request) {
	now := h.clock.Now()
	bookings := h.store.ListBookings()

	upcoming := make([]models.Booking, 0, len(bookings))
	for _, b := range bookings {
		if b.StartTime.After(now) {
			upcoming = append(upcoming, b)
		}
	}

	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].StartTime.Before(upcoming[j].StartTime)
	})

	writeJSON(w, http.StatusOK, upcoming)
}

// Delete cancels a booking by id.
func (h *AdminBookingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id_required", "id is required")
		return
	}

	if !h.store.DeleteBooking(id) {
		writeError(w, http.StatusNotFound, "not_found", "booking not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
