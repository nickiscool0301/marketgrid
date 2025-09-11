package http

import (
	"marketgrid/user/app/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminApp *admin.AdminApp
}

func NewAdminHandler(adminApp *admin.AdminApp) *AdminHandler {
	return &AdminHandler{
		adminApp: adminApp,
	}
}

func (h *AdminHandler) RegisterRoutes(router *gin.Engine) {
	adminV1 := router.Group("/api/v1/admin")
	{
		adminV1.GET("/users", h.GetAllUsersForAdmin)
	}
}

func (h *AdminHandler) GetAllUsersForAdmin(c *gin.Context) {
	users, err := h.adminApp.GetAllUsersForAdmin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
