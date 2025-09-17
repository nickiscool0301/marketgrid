package postgres

import (
	"context"
	"errors"
	"time"

	"marketgrid/user/domain/model"
	"marketgrid/user/domain/port"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// gormUser is the persistence model used only in this layer
// It contains GORM-specific tags and mirrors the database schema.
type gormUser struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id"`
	Email     string    `gorm:"type:text;uniqueIndex;not null;column:email"`
	Password  string    `gorm:"type:text;not null;column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (gormUser) TableName() string { return "users" }

// PostgresUserRepository implements the UserRepository port using PostgreSQL
type PostgresUserRepository struct {
	db *gorm.DB
}

// Ensure PostgresUserRepository implements the UserRepository interface
var _ port.UserRepository = (*PostgresUserRepository)(nil)

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

func (r *PostgresUserRepository) GetAllEmails(ctx context.Context) ([]string, error) {
	var emails []string
	if err := r.db.WithContext(ctx).Model(&gormUser{}).Pluck("email", &emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

// Helper functions for mapping between domain and persistence models
func toGorm(u *model.User) *gormUser {
	return &gormUser{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreateAt,
		UpdatedAt: u.UpdateAt,
	}
}

func toDomain(u *gormUser) *model.User {
	return &model.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		CreateAt: u.CreatedAt,
		UpdateAt: u.UpdatedAt,
	}
}
