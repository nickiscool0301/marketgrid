package user

import (
	"context"
	"errors"
	"marketgrid/user/internal/app/dto"
	"marketgrid/user/internal/domain/model"
	"marketgrid/user/internal/domain/port"
)

type UserApp struct {
	userRepo              port.UserRepository
	emailExistenceService port.EmailExistenceService
}

func NewUserApp(userRepo port.UserRepository, emailExistenceService port.EmailExistenceService) *UserApp {
	return &UserApp{
		userRepo:              userRepo,
		emailExistenceService: emailExistenceService,
	}
}

func (app *UserApp) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := app.emailExistenceService.EmailExists(ctx, req.Email)
	if err != nil {
		if _, dbErr := app.userRepo.FindByEmail(ctx, req.Email); dbErr == nil {
			return nil, errors.New("user with this email already exists")
		}
	} else if exists {
		return nil, errors.New("user with this email already exists")
	}

	newUser, err := model.NewUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	if err := app.userRepo.SaveUser(ctx, newUser); err != nil {
		return nil, err
	}

	// Add email to existence service after successful creation
	if addErr := app.emailExistenceService.AddEmail(ctx, req.Email); addErr != nil {
		// Log error but don't fail the request since user was created successfully
		// In production, you might want to use a proper logger here
	}

	return &dto.UserResponse{
		ID:       newUser.ID,
		Email:    newUser.Email,
		CreateAt: newUser.CreateAt,
		UpdateAt: newUser.UpdateAt,
	}, nil
}

func (app *UserApp) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := app.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}

func (app *UserApp) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := app.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}, nil
}
