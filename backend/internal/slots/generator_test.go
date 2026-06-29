package slots_test

import (
	"testing"
	"time"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/clock"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/slots"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func startTimes(slots []slots.Slot) []string {
	out := make([]string, len(slots))
	for i, s := range slots {
		out[i] = s.StartTime.UTC().Format(time.RFC3339)
	}
	return out
}

func statuses(slots []slots.Slot) []string {
	out := make([]string, len(slots))
	for i, s := range slots {
		out[i] = s.Status
	}
	return out
}

func TestGenerate(t *testing.T) {
	moscow := mustLocation("Europe/Moscow")

	tests := []struct {
		name     string
		now      time.Time
		duration time.Duration
		taken    []time.Time
		wantDays int
		check    func(t *testing.T, days []slots.AvailableDay)
	}{
		{
			name:     "14 calendar days from Monday with working day slots",
			now:      time.Date(2026, 6, 29, 8, 0, 0, 0, moscow),
			duration: 60 * time.Minute,
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				require.Len(t, days, 14)
				assert.Equal(t, "2026-06-29", days[0].Date)
				assert.Equal(t, "2026-07-12", days[13].Date)

				// Monday 2026-06-29: 09:00-17:00 MSK -> 06:00-14:00 UTC, 9 slots.
				require.Len(t, days[0].Slots, 9)
				assert.Equal(t, "2026-06-29T06:00:00Z", days[0].Slots[0].StartTime.UTC().Format(time.RFC3339))
				assert.Equal(t, "free", days[0].Slots[0].Status)
				assert.Equal(t, "2026-06-29T14:00:00Z", days[0].Slots[8].StartTime.UTC().Format(time.RFC3339))
				assert.Equal(t, "free", days[0].Slots[8].Status)

				// Saturday 2026-07-04 is empty.
				assert.Equal(t, "2026-07-04", days[5].Date)
				assert.Empty(t, days[5].Slots)

				// Sunday 2026-07-05 is empty.
				assert.Equal(t, "2026-07-05", days[6].Date)
				assert.Empty(t, days[6].Slots)
			},
		},
		{
			name:     "last slot for 60-minute event starts at 17:00",
			now:      time.Date(2026, 6, 29, 8, 0, 0, 0, moscow),
			duration: 60 * time.Minute,
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				require.NotEmpty(t, days[0].Slots)
				last := days[0].Slots[len(days[0].Slots)-1]
				assert.Equal(t, "2026-06-29T14:00:00Z", last.StartTime.UTC().Format(time.RFC3339))
				assert.Equal(t, "free", last.Status)
			},
		},
		{
			name:     "last slot for 9-hour event starts at 09:00 only",
			now:      time.Date(2026, 6, 29, 8, 0, 0, 0, moscow),
			duration: 540 * time.Minute,
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				// Working day has exactly one slot starting at 09:00 MSK (06:00 UTC).
				require.Len(t, days[0].Slots, 1)
				assert.Equal(t, "2026-06-29T06:00:00Z", days[0].Slots[0].StartTime.UTC().Format(time.RFC3339))
				assert.Equal(t, "free", days[0].Slots[0].Status)
			},
		},
		{
			name:     "past slots for today are excluded",
			now:      time.Date(2026, 6, 29, 10, 30, 0, 0, moscow),
			duration: 60 * time.Minute,
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				// 09:00 and 10:00 MSK are in the past; 11:00 MSK is first.
				require.Len(t, days[0].Slots, 7)
				assert.Equal(t, "2026-06-29T08:00:00Z", days[0].Slots[0].StartTime.UTC().Format(time.RFC3339))
				assert.Equal(t, "free", days[0].Slots[0].Status)
			},
		},
		{
			name:     "taken slots are marked as taken",
			now:      time.Date(2026, 6, 29, 8, 0, 0, 0, moscow),
			duration: 60 * time.Minute,
			taken: []time.Time{
				time.Date(2026, 6, 29, 11, 0, 0, 0, moscow), // 11:00 MSK = 08:00 UTC.
			},
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				// Total 9 slots remain generated; the 11:00 MSK one is taken.
				require.Len(t, days[0].Slots, 9)
				starts := startTimes(days[0].Slots)
				assert.Contains(t, starts, "2026-06-29T08:00:00Z")
				assert.Contains(t, starts, "2026-06-29T06:00:00Z")
				assert.Contains(t, starts, "2026-06-29T14:00:00Z")

				statuses := statuses(days[0].Slots)
				assert.Contains(t, statuses, "taken")
				assert.Contains(t, statuses, "free")

				for _, slot := range days[0].Slots {
					expected := "free"
					if slot.StartTime.UTC().Equal(time.Date(2026, 6, 29, 8, 0, 0, 0, time.UTC)) {
						expected = "taken"
					}
					assert.Equal(t, expected, slot.Status, "slot %s", slot.StartTime.UTC().Format(time.RFC3339))
				}
			},
		},
		{
			name:     "fully booked working day has all slots marked as taken",
			now:      time.Date(2026, 6, 29, 8, 0, 0, 0, moscow),
			duration: 60 * time.Minute,
			taken: []time.Time{
				time.Date(2026, 6, 29, 9, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 10, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 11, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 12, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 13, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 14, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 15, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 16, 0, 0, 0, moscow),
				time.Date(2026, 6, 29, 17, 0, 0, 0, moscow),
			},
			wantDays: 14,
			check: func(t *testing.T, days []slots.AvailableDay) {
				assert.Equal(t, "2026-06-29", days[0].Date)
				require.Len(t, days[0].Slots, 9)
				for _, slot := range days[0].Slots {
					assert.Equal(t, "taken", slot.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := clock.NewMockClock(tt.now)
			days := slots.Generate("et-1", tt.duration, moscow, mc, tt.taken)
			require.Len(t, days, tt.wantDays)
			tt.check(t, days)
		})
	}
}
