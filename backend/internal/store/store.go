package store

import (
	"errors"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
)

// Common store errors.
var (
	ErrNotFound  = errors.New("not found")
	ErrSlotTaken = errors.New("slot is already taken")
)

// Store defines the contract for persistence. Implementations must be safe for
// concurrent use.
type Store interface {
	ListEventTypes() []models.EventType
	GetEventType(id string) (models.EventType, bool)
	CreateEventType(et models.EventType)
	UpdateEventType(et models.EventType) bool
	DeleteEventType(id string) bool

	ListBookings() []models.Booking
	GetBooking(id string) (models.Booking, bool)
	CreateBooking(b models.Booking) error
	DeleteBooking(id string) bool
}
