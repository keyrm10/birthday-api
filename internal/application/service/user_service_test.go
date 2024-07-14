package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username user.Username) (user.User, error) {
	args := m.Called(username)
	return args.Get(0).(user.User), args.Error(1)
}

func TestUserService(t *testing.T) {
	username, _ := user.NewUsername("john")
	now := time.Now()

	t.Run("SaveUser", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		dateOfBirth := now.AddDate(-30, 1, 1).Format(time.DateOnly)

		mockRepo.On("Save", mock.AnythingOfType("user.User")).Return(nil)

		err := service.SaveUser(username, dateOfBirth)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetBirthdayMessage_Today", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		dob, _ := user.NewDateOfBirth(now.AddDate(-30, 0, 0).Format(time.DateOnly))
		mockUser := user.User{Username: username, DateOfBirth: dob}

		mockRepo.On("FindByUsername", username).Return(mockUser, nil)

		message, err := service.GetBirthdayMessage(username)

		assert.NoError(t, err)
		assert.Contains(t, message, "Happy birthday!")
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetBirthdayMessage_Future", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		dob, _ := user.NewDateOfBirth(now.AddDate(-30, 0, 5).Format(time.DateOnly))
		mockUser := user.User{Username: username, DateOfBirth: dob}

		mockRepo.On("FindByUsername", username).Return(mockUser, nil)

		message, err := service.GetBirthdayMessage(username)

		assert.NoError(t, err)
		assert.Contains(t, message, "Your birthday is in")
		mockRepo.AssertExpectations(t)
	})
}
