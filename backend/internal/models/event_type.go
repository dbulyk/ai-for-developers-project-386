package models

import "fmt"

// EventType describes a bookable service and its duration in minutes.
type EventType struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int32  `json:"durationMinutes"`
}

// NewEventType validates and creates a new EventType.
func NewEventType(id, name, description string, duration int32) (EventType, error) {
	et := EventType{
		ID:              id,
		Name:            name,
		Description:     description,
		DurationMinutes: duration,
	}
	if err := et.validate(); err != nil {
		return EventType{}, err
	}
	return et, nil
}

// Update modifies the event type in place and re-validates it.
func (e *EventType) Update(name, description string, duration int32) error {
	e.Name = name
	e.Description = description
	e.DurationMinutes = duration
	return e.validate()
}

func (e *EventType) validate() error {
	if e.ID == "" {
		return NewValidationError("id_required", "id is required")
	}
	if e.Name == "" {
		return NewValidationError("name_required", "name is required")
	}
	if e.DurationMinutes <= 0 {
		return NewValidationError("duration_invalid", fmt.Sprintf("duration must be positive, got %d", e.DurationMinutes))
	}
	return nil
}
