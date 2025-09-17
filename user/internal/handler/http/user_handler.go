package http

import (
	"marketgrid/user/internal/app/dto"
	"marketgrid/user/internal/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userApp *user.UserApp
}

func NewUserHandler(userApp *user.UserApp) *UserHandler {
	return &UserHandler{
		userApp: userApp,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1/user")
	{
		v1.POST("/users", h.CreateUser)
		v1.GET("/users/:id", h.GetUserByID)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userApp.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}
	user, err := h.userApp.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
