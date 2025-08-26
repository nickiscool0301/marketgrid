package persistence

import (
	"time"

	"marketgrid/user/internal/domain/model"
)

// gormUser is the persistence model used only in this adapter layer
// It contains GORM-specific tags and mirrors the database schema.
type gormUser struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id"`
	Email     string    `gorm:"type:text;uniqueIndex;not null;column:email"`
	Password  string    `gorm:"type:text;not null;column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (gormUser) TableName() string { return "users" }

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
