package user

import (
	"errors"
	"regexp"
	"time"
)

type Username string

func NewUsername(name string) (Username, error) {
	if !isValidUsername(name) {
		return "", errors.New("invalid username")
	}
	return Username(name), nil
}

func isValidUsername(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(username)
}

type DateOfBirth struct {
	time.Time
}

func NewDateOfBirth(date string) (DateOfBirth, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		t, err = time.Parse(time.RFC3339, date)
		if err != nil {
			return DateOfBirth{}, err
		}
	}
	if t.After(time.Now()) {
		return DateOfBirth{}, errors.New("date of birth must be in the past")
	}
	return DateOfBirth{Time: t}, nil
}

type User struct {
	Username    Username
	DateOfBirth DateOfBirth
}

func (u User) DaysUntilBirthday() int {
	today := time.Now().Truncate(time.Hour * 24)
	nextBirthday := time.Date(today.Year(), u.DateOfBirth.Month(), u.DateOfBirth.Day(), 0, 0, 0, 0, time.UTC)
	if nextBirthday.Before(today) {
		nextBirthday = nextBirthday.AddDate(1, 0, 0)
	}
	return int(nextBirthday.Sub(today).Hours() / 24)
}

func (u User) IsBirthdayToday() bool {
	now := time.Now()
	return now.Month() == u.DateOfBirth.Month() && now.Day() == u.DateOfBirth.Day()
}
