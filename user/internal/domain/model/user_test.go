package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Valid(t *testing.T) {
	u, err := NewUser("test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", u.Email)
	assert.NotEmpty(t, u.Password)
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := NewUser("bad-email", "password")
	assert.Error(t, err)
}

func TestAuthenticate(t *testing.T) {
	u, err := NewUser("test@example.com", "secret")
	require.NoError(t, err)
	assert.NoError(t, u.Authenticate("secret"))
	assert.Error(t, u.Authenticate("wrong"))
}
