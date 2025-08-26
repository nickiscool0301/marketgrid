package model

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Email    string
	Password string
	CreateAt time.Time
	UpdateAt time.Time
}

func NewUser(email, plainTextPassword string) (*User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}
	hashedPass, err := hashPassword(plainTextPassword)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: hashedPass,
		CreateAt: time.Now().UTC(),
		UpdateAt: time.Now().UTC(),
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (user *User) Authenticate(plainTextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainTextPassword))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}
