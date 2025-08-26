package port

import (
	"context"
	"marketgrid/user/internal/domain/model"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindAll(ctx context.Context) ([]*model.User, error)
}
