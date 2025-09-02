package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"marketgrid/user/internal/application/dto"
	"marketgrid/user/internal/application/port"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeService struct{}

func (f *fakeService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	return &dto.UserResponse{ID: "1", Email: req.Email}, nil
}
func (f *fakeService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	return []*dto.UserResponse{{ID: "1", Email: "a@b.com"}}, nil
}
func (f *fakeService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	return &dto.UserResponse{ID: id, Email: "a@b.com"}, nil
}
func (f *fakeService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	return &dto.UserResponse{ID: "1", Email: email}, nil
}

var _ port.UserService = (*fakeService)(nil)

func TestRegisterRoutes_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(&fakeService{})
	h.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "a@b.com")
}
