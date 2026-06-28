package slots

import (
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
)

// AvailableDay groups free slot start times for a single calendar day.
type AvailableDay struct {
	Date  string   `json:"date"`
	Slots []string `json:"slots"`
}

// Generate returns available booking slots for the next 14 calendar days
// (including today), grouped by calendar day in the provided timezone.
// Slots are generated on working days (Mon-Fri) from 09:00 to 18:00 local
// time with the given event duration as the grid step. Slots that have
// already passed today, or are present in taken, are omitted.
func Generate(eventTypeID string, eventTypeDuration time.Duration, tz *time.Location, clock clock.Clock, taken []time.Time) []AvailableDay {
	_ = eventTypeID

	now := clock.Now()
	today := now.In(tz)

	const (
		workStartHour = 9
		workEndHour   = 18
	)
	workDuration := time.Duration(workEndHour-workStartHour) * time.Hour

	days := make([]AvailableDay, 0, 14)

	for offset := 0; offset < 14; offset++ {
		date := today.AddDate(0, 0, offset)

		day := AvailableDay{
			Date:  date.Format(time.DateOnly),
			Slots: make([]string, 0),
		}

		if eventTypeDuration > 0 && eventTypeDuration <= workDuration && date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			dayStart := time.Date(date.Year(), date.Month(), date.Day(), workStartHour, 0, 0, 0, tz)
			dayEnd := time.Date(date.Year(), date.Month(), date.Day(), workEndHour, 0, 0, 0, tz)
			lastStart := dayEnd.Add(-eventTypeDuration)

			for slot := dayStart; !slot.After(lastStart); slot = slot.Add(eventTypeDuration) {
				if isToday(date, today) && (slot.Before(now) || slot.Equal(now)) {
					continue
				}
				if isTaken(slot, taken) {
					continue
				}
				day.Slots = append(day.Slots, slot.UTC().Format(time.RFC3339))
			}
		}

		days = append(days, day)
	}

	return days
}

func isToday(date, today time.Time) bool {
	return date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()
}

func isTaken(slot time.Time, taken []time.Time) bool {
	for _, t := range taken {
		if t.Equal(slot) {
			return true
		}
	}
	return false
}
