package admin

import (
	"context"
	"marketgrid/user/internal/app/dto"
	"marketgrid/user/internal/domain/port"
)

type AdminApp struct {
	userRepo port.UserRepository
}

func NewAdminApp(userRepo port.UserRepository) *AdminApp {
	return &AdminApp{
		userRepo: userRepo,
	}
}

func (app *AdminApp) GetAllUsersForAdmin(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := app.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			CreateAt: user.CreateAt,
			UpdateAt: user.UpdateAt,
		})
	}

	return userResponses, nil
}
