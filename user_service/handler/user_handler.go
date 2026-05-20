package handler

import (
	"net/http"
	"user_service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(u service.UserService) *UserHandler {
	return &UserHandler{userService: u}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetByID(id) // <--- Di sini letak error sebelumnya
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, user)
}
