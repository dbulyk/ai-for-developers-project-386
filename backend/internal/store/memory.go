package store

import (
	"sync"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
)

// MemoryStore is an in-memory implementation of Store protected by a mutex.
type MemoryStore struct {
	mu         sync.RWMutex
	eventTypes map[string]models.EventType
	bookings   map[string]models.Booking
}

// NewMemoryStore creates a new empty MemoryStore.
func NewMemoryStore() Store {
	return &MemoryStore{
		eventTypes: make(map[string]models.EventType),
		bookings:   make(map[string]models.Booking),
	}
}

// ListEventTypes returns all stored event types.
func (s *MemoryStore) ListEventTypes() []models.EventType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]models.EventType, 0, len(s.eventTypes))
	for _, et := range s.eventTypes {
		list = append(list, et)
	}
	return list
}

// GetEventType returns an event type by id.
func (s *MemoryStore) GetEventType(id string) (models.EventType, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	et, ok := s.eventTypes[id]
	return et, ok
}

// CreateEventType stores a new event type.
func (s *MemoryStore) CreateEventType(et models.EventType) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventTypes[et.ID] = et
}

// UpdateEventType replaces an existing event type. Returns false if it does not exist.
func (s *MemoryStore) UpdateEventType(et models.EventType) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventTypes[et.ID]; !ok {
		return false
	}
	s.eventTypes[et.ID] = et
	return true
}

// DeleteEventType removes an event type by id. Returns false if it does not exist.
func (s *MemoryStore) DeleteEventType(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventTypes[id]; !ok {
		return false
	}
	delete(s.eventTypes, id)
	return true
}

// ListBookings returns all stored bookings.
func (s *MemoryStore) ListBookings() []models.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]models.Booking, 0, len(s.bookings))
	for _, b := range s.bookings {
		list = append(list, b)
	}
	return list
}

// GetBooking returns a booking by id.
func (s *MemoryStore) GetBooking(id string) (models.Booking, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.bookings[id]
	return b, ok
}

// CreateBooking stores a booking if its start time is not already taken.
// Returns ErrSlotTaken when another booking occupies the same slot.
func (s *MemoryStore) CreateBooking(b models.Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.bookings {
		if existing.StartTime.Equal(b.StartTime) {
			return ErrSlotTaken
		}
	}
	s.bookings[b.ID] = b
	return nil
}

// DeleteBooking removes a booking by id. Returns false if it does not exist.
func (s *MemoryStore) DeleteBooking(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.bookings[id]; !ok {
		return false
	}
	delete(s.bookings, id)
	return true
}
