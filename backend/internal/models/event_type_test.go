package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventType(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		n           string
		description string
		duration    int32
		wantErr     bool
		errCode     string
	}{
		{
			name:        "valid",
			id:          "et-1",
			n:           "Meeting",
			description: "Quick meeting",
			duration:    30,
			wantErr:     false,
		},
		{
			name:     "missing id",
			id:       "",
			n:        "Meeting",
			duration: 30,
			wantErr:  true,
			errCode:  "id_required",
		},
		{
			name:     "missing name",
			id:       "et-1",
			n:        "",
			duration: 30,
			wantErr:  true,
			errCode:  "name_required",
		},
		{
			name:     "zero duration",
			id:       "et-1",
			n:        "Meeting",
			duration: 0,
			wantErr:  true,
			errCode:  "duration_invalid",
		},
		{
			name:     "negative duration",
			id:       "et-1",
			n:        "Meeting",
			duration: -10,
			wantErr:  true,
			errCode:  "duration_invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et, err := NewEventType(tt.id, tt.n, tt.description, tt.duration)
			if tt.wantErr {
				require.Error(t, err)
				var vErr *ValidationError
				require.ErrorAs(t, err, &vErr)
				assert.Equal(t, tt.errCode, vErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.id, et.ID)
			assert.Equal(t, tt.n, et.Name)
			assert.Equal(t, tt.description, et.Description)
			assert.Equal(t, tt.duration, et.DurationMinutes)
		})
	}
}

func TestEventType_Update(t *testing.T) {
	et, err := NewEventType("et-1", "Meeting", "Description", 30)
	require.NoError(t, err)

	err = et.Update("Updated", "New description", 60)
	require.NoError(t, err)
	assert.Equal(t, "Updated", et.Name)
	assert.Equal(t, "New description", et.Description)
	assert.Equal(t, int32(60), et.DurationMinutes)

	err = et.Update("", "Desc", 30)
	require.Error(t, err)
	var vErr *ValidationError
	require.ErrorAs(t, err, &vErr)
	assert.Equal(t, "name_required", vErr.Code)
}
