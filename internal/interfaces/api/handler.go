package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

type UserService interface {
	SaveUser(username user.Username, dateOfBirth string) error
	GetBirthdayMessage(username user.Username) (string, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{userService: userService}
}

type saveUserRequest struct {
	DateOfBirth string `json:"dateOfBirth"`
}

func (h *Handler) SaveUser(c *gin.Context) {
	username, err := user.NewUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req saveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.SaveUser(username, req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetUserBirthday(c *gin.Context) {
	username, err := user.NewUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.userService.GetBirthdayMessage(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
