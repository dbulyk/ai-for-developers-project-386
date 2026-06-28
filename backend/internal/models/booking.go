package models

import (
	"fmt"
	"time"
)

// Booking represents a confirmed appointment.
type Booking struct {
	ID          string
	EventTypeID string
	GuestName   string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
}

// NewBooking validates and creates a new Booking. The end time is computed
// from startTime and duration; CreatedAt is set to now.
func NewBooking(id, eventTypeID, guestName string, startTime time.Time, duration time.Duration, now time.Time) (Booking, error) {
	b := Booking{
		ID:          id,
		EventTypeID: eventTypeID,
		GuestName:   guestName,
		StartTime:   startTime,
		EndTime:     startTime.Add(duration),
		CreatedAt:   now,
	}
	if err := b.validate(now, duration); err != nil {
		return Booking{}, err
	}
	return b, nil
}

func (b *Booking) validate(now time.Time, duration time.Duration) error {
	if b.ID == "" {
		return NewValidationError("id_required", "id is required")
	}
	if b.EventTypeID == "" {
		return NewValidationError("event_type_id_required", "event type id is required")
	}
	if b.GuestName == "" {
		return NewValidationError("guest_name_required", "guest name is required")
	}
	if duration <= 0 {
		return NewValidationError("duration_invalid", fmt.Sprintf("duration must be positive, got %s", duration))
	}
	if !b.StartTime.After(now) {
		return NewValidationError("start_time_invalid", "start time must be in the future")
	}
	return nil
}
