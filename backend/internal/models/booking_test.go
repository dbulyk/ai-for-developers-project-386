package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBooking(t *testing.T) {
	now := time.Date(2026, 6, 28, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		id          string
		eventTypeID string
		guestName   string
		start       time.Time
		duration    time.Duration
		wantErr     bool
		errCode     string
	}{
		{
			name:        "valid",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now.Add(time.Hour),
			duration:    30 * time.Minute,
			wantErr:     false,
		},
		{
			name:        "missing id",
			id:          "",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now.Add(time.Hour),
			duration:    30 * time.Minute,
			wantErr:     true,
			errCode:     "id_required",
		},
		{
			name:        "missing event type id",
			id:          "b-1",
			eventTypeID: "",
			guestName:   "Alice",
			start:       now.Add(time.Hour),
			duration:    30 * time.Minute,
			wantErr:     true,
			errCode:     "event_type_id_required",
		},
		{
			name:        "missing guest name",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "",
			start:       now.Add(time.Hour),
			duration:    30 * time.Minute,
			wantErr:     true,
			errCode:     "guest_name_required",
		},
		{
			name:        "start in the past",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now.Add(-time.Hour),
			duration:    30 * time.Minute,
			wantErr:     true,
			errCode:     "start_time_invalid",
		},
		{
			name:        "start equal to now",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now,
			duration:    30 * time.Minute,
			wantErr:     true,
			errCode:     "start_time_invalid",
		},
		{
			name:        "zero duration",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now.Add(time.Hour),
			duration:    0,
			wantErr:     true,
			errCode:     "duration_invalid",
		},
		{
			name:        "negative duration",
			id:          "b-1",
			eventTypeID: "et-1",
			guestName:   "Alice",
			start:       now.Add(time.Hour),
			duration:    -time.Minute,
			wantErr:     true,
			errCode:     "duration_invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBooking(tt.id, tt.eventTypeID, tt.guestName, tt.start, tt.duration, now)
			if tt.wantErr {
				require.Error(t, err)
				var vErr *ValidationError
				require.ErrorAs(t, err, &vErr)
				assert.Equal(t, tt.errCode, vErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.id, b.ID)
			assert.Equal(t, tt.eventTypeID, b.EventTypeID)
			assert.Equal(t, tt.guestName, b.GuestName)
			assert.Equal(t, tt.start, b.StartTime)
			assert.Equal(t, tt.start.Add(tt.duration), b.EndTime)
			assert.WithinDuration(t, now, b.CreatedAt, time.Second)
		})
	}
}
