package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SaveUser(username user.Username, dateOfBirth string) error {
	args := m.Called(username, dateOfBirth)
	return args.Error(0)
}

func (m *MockUserService) GetBirthdayMessage(username user.Username) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func TestHandler(t *testing.T) {
	t.Run("SaveUser", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewHandler(mockService)

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.PUT("/hello/:username", handler.SaveUser)

		mockService.On("SaveUser", user.Username("john"), "1990-01-01").Return(nil)

		body := saveUserRequest{DateOfBirth: "1990-01-01"}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPut, "/hello/john", strings.NewReader(string(jsonBody)))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GetUserBirthday", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewHandler(mockService)

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/hello/:username", handler.GetUserBirthday)

		mockService.On("GetBirthdayMessage", user.Username("john")).Return("Hello, john! Your birthday is in 5 days", nil)

		req, _ := http.NewRequest(http.MethodGet, "/hello/john", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Hello, john!")
		assert.True(t, strings.Contains(response["message"], "Happy birthday!") || strings.Contains(response["message"], "Your birthday is in"))
		mockService.AssertExpectations(t)
	})
}
