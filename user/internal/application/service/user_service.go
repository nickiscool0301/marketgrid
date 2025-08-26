package service

import (
	"context"
	"errors"
	"marketgrid/user/internal/application/dto"
	"marketgrid/user/internal/application/port"
	"marketgrid/user/internal/domain/model"
)

type userService struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	newUser, err := model.NewUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.SaveUser(ctx, newUser); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       newUser.ID,
		Email:    newUser.Email,
		CreateAt: newUser.CreateAt,
		UpdateAt: newUser.UpdateAt,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
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

func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
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

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
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
