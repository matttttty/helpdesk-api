package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"helpdesk-api/internal/model"
	"helpdesk-api/internal/repository"
)

func TestUserRepository_CreateAndGetByEmail(t *testing.T) {
	repo := repository.NewUserRepository(testDB) // ⚠ перевір назву конструктора у себе
	ctx := context.Background()

	user := &model.User{
		Name:     "Alice",
		Email:    "alice@test.com",
		Password: "hashed-password",
		Role:     model.RoleClient,
	}

	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	// SELECT назад
	got, err := repo.GetUserByEmail(ctx, "alice@test.com")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "alice@test.com", got.Email)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, model.RoleClient, got.Role)
	assert.NotZero(t, got.ID)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	repo := repository.NewUserRepository(testDB)

	got, err := repo.GetUserByEmail(context.Background(), "ghost@test.com")

	require.NoError(t, err)
	assert.Nil(t, got)
}
