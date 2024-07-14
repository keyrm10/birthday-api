package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("NewUsername", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			wantErr  bool
			expected Username
		}{
			{"Valid username", "john", false, "john"},
			{"Invalid username", "john!123", true, ""},
			{"Empty username", "", true, ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewUsername(tt.input)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, got)
				}
			})
		}
	})

	t.Run("NewDateOfBirth", func(t *testing.T) {
		tests := []struct {
			name    string
			input   string
			wantErr bool
		}{
			{"Valid date", "1990-01-01", false},
			{"Invalid date format", "1990/01/01", true},
			{"Future date", "2100-01-01", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := NewDateOfBirth(tt.input)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("DaysUntilBirthday", func(t *testing.T) {
		now := time.Now()

		tomorrowUser := User{
			Username:    "tomorrow",
			DateOfBirth: DateOfBirth{Time: time.Date(1990, now.Month(), now.AddDate(0, 0, 1).Day(), 0, 0, 0, 0, time.UTC)},
		}
		assert.Equal(t, 1, tomorrowUser.DaysUntilBirthday())

		todayUser := User{
			Username:    "today",
			DateOfBirth: DateOfBirth{Time: time.Date(1990, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)},
		}
		assert.Equal(t, 0, todayUser.DaysUntilBirthday())

		yesterdayUser := User{
			Username:    "yesterday",
			DateOfBirth: DateOfBirth{Time: time.Date(1990, now.Month(), now.AddDate(0, 0, -1).Day(), 0, 0, 0, 0, time.UTC)},
		}
		days := yesterdayUser.DaysUntilBirthday()
		// Account for leap years
		assert.True(t, days == 364 || days == 365)
	})

	t.Run("IsBirthdayToday", func(t *testing.T) {
		now := time.Now()

		u := User{
			Username:    "jane",
			DateOfBirth: DateOfBirth{Time: time.Date(now.Year()-30, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)},
		}

		assert.True(t, u.IsBirthdayToday())
	})
}
