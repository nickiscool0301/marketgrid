package persistence

import (
	"context"
	"errors"

	"marketgrid/user/internal/domain/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(dataSourceName string) (*PostgresUserRepository, error) {
	dial := postgres.Open(dataSourceName)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgresUserRepository{db: db}, nil
}

func (r *PostgresUserRepository) Init() error {
	return r.db.AutoMigrate(&gormUser{})
}

func (r *PostgresUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(toGorm(user)).Error
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u gormUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toDomain(&u), nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u gormUser
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toDomain(&u), nil
}

func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	var us []gormUser
	if err := r.db.WithContext(ctx).Find(&us).Error; err != nil {
		return nil, err
	}
	result := make([]*model.User, 0, len(us))
	for i := range us {
		result = append(result, toDomain(&us[i]))
	}
	return result, nil
}
