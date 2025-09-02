package service

import (
	"context"
	"errors"
	"testing"

	"marketgrid/user/internal/application/dto"
	"marketgrid/user/internal/application/port"
	"marketgrid/user/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	users map[string]*model.User
}

func (f *fakeRepo) SaveUser(ctx context.Context, u *model.User) error {
	if f.users == nil {
		f.users = map[string]*model.User{}
	}
	f.users[u.ID] = u
	return nil
}
func (f *fakeRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}
func (f *fakeRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeRepo) FindAll(ctx context.Context) ([]*model.User, error) {
	var res []*model.User
	for _, u := range f.users {
		res = append(res, u)
	}
	return res, nil
}

var _ port.UserRepository = (*fakeRepo)(nil)

func TestCreateUser_Success(t *testing.T) {
	svc := NewUserService(&fakeRepo{})
	res, err := svc.CreateUser(context.Background(), &dto.CreateUserRequest{Email: "a@b.com", Password: "x"})
	require.NoError(t, err)
	assert.Equal(t, "a@b.com", res.Email)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo := &fakeRepo{}
	u, _ := model.NewUser("dup@b.com", "p")
	_ = repo.SaveUser(context.Background(), u)
	svc := NewUserService(repo)
	_, err := svc.CreateUser(context.Background(), &dto.CreateUserRequest{Email: "dup@b.com", Password: "x"})
	assert.Error(t, err)
}
