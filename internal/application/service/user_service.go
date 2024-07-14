package service

import (
	"fmt"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

type UserRepository interface {
	Save(user user.User) error
	FindByUsername(username user.Username) (user.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SaveUser(username user.Username, dateOfBirth string) error {
	dob, err := user.NewDateOfBirth(dateOfBirth)
	if err != nil {
		return err
	}

	u := user.User{Username: username, DateOfBirth: dob}
	return s.repo.Save(u)
}

func (s *UserService) GetBirthdayMessage(username user.Username) (string, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	if u.IsBirthdayToday() {
		return fmt.Sprintf("Hello, %s! Happy birthday!", u.Username), nil
	}

	daysUntil := u.DaysUntilBirthday()
	return fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", u.Username, daysUntil), nil
}
