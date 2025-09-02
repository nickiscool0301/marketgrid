package persistence

import (
	"context"
	"testing"

	"marketgrid/user/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&gormUser{}))
	return db
}

func TestSaveAndFindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := &PostgresUserRepository{db: db}
	u, _ := model.NewUser("t@e.com", "p")
	require.NoError(t, repo.SaveUser(context.Background(), u))
	got, err := repo.FindByID(context.Background(), u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.Email, got.Email)
}

func TestSaveAndFindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := &PostgresUserRepository{db: db}
	u, _ := model.NewUser("t@e.com", "p")
	require.NoError(t, repo.SaveUser(context.Background(), u))
	got, err := repo.FindByEmail(context.Background(), u.Email)
	require.NoError(t, err)
	assert.Equal(t, u.Email, got.Email)
}

func TestFindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := &PostgresUserRepository{db: db}
	u1, _ := model.NewUser("a@e.com", "p1")
	u2, _ := model.NewUser("b@e.com", "p2")
	require.NoError(t, repo.SaveUser(context.Background(), u1))
	require.NoError(t, repo.SaveUser(context.Background(), u2))
	s, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, s, 2)
	emails := map[string]bool{s[0].Email: true, s[1].Email: true}
	assert.True(t, emails["a@e.com"])
	assert.True(t, emails["b@e.com"])
}
