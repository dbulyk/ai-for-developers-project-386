package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/slots"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// PublicBookingsHandler handles public booking creation.
type PublicBookingsHandler struct {
	store store.Store
	clock clock.Clock
	tz    *time.Location
}

// NewPublicBookingsHandler creates a new PublicBookingsHandler.
func NewPublicBookingsHandler(s store.Store, c clock.Clock, tz *time.Location) *PublicBookingsHandler {
	return &PublicBookingsHandler{
		store: s,
		clock: c,
		tz:    tz,
	}
}

// RegisterRoutes mounts the public bookings routes on the given router.
func (h *PublicBookingsHandler) RegisterRoutes(r chi.Router) {
	r.Post("/public/bookings", h.Create)
}

// createBookingRequest is the request body for creating a booking.
type createBookingRequest struct {
	EventTypeID string `json:"eventTypeId"`
	GuestName   string `json:"guestName"`
	StartTime   string `json:"startTime"`
}

// Create creates a new booking for a public slot.
func (h *PublicBookingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	et, ok := h.store.GetEventType(req.EventTypeID)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "event type not found")
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "start_time_invalid", "start time must be a valid RFC3339 timestamp")
		return
	}

	valid, taken := h.isValidSlot(et, startTime)
	if taken {
		writeError(w, http.StatusConflict, "slot_taken", "slot is already taken")
		return
	}
	if !valid {
		writeError(w, http.StatusBadRequest, "slot_invalid", "requested slot is not available")
		return
	}

	id, err := generateID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to generate id")
		return
	}

	now := h.clock.Now()
	duration := time.Duration(et.DurationMinutes) * time.Minute
	booking, err := models.NewBooking(id, et.ID, req.GuestName, startTime.UTC(), duration, now)
	if err != nil {
		var vErr *models.ValidationError
		if errors.As(err, &vErr) {
			writeError(w, http.StatusBadRequest, vErr.Code, vErr.Message)
			return
		}
		writeError(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	if err := h.store.CreateBooking(booking); err != nil {
		if errors.Is(err, store.ErrSlotTaken) {
			writeError(w, http.StatusConflict, "slot_taken", "slot is already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to create booking")
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

// isValidSlot verifies that startTime is one of the currently generated
// slots for the given event type. It returns (true, false) when the slot is
// free, (false, true) when the slot is already taken, and (false, false)
// when it is outside the available window.
func (h *PublicBookingsHandler) isValidSlot(et models.EventType, startTime time.Time) (bool, bool) {
	startTime = startTime.UTC()
	bookings := h.store.ListBookings()
	taken := make([]time.Time, 0, len(bookings))
	for _, b := range bookings {
		if b.EventTypeID == et.ID {
			taken = append(taken, b.StartTime)
		}
	}

	duration := time.Duration(et.DurationMinutes) * time.Minute
	available := slots.Generate(et.ID, duration, h.tz, h.clock, taken)
	for _, day := range available {
		for _, slot := range day.Slots {
			if slot.StartTime.UTC().Equal(startTime) {
				return slot.Status == "free", slot.Status == "taken"
			}
		}
	}

	return false, false
}
