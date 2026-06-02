package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"helpdesk-api/internal/model"
	"helpdesk-api/internal/repository"
	"helpdesk-api/internal/service"
)

func TestTicketRepository_CreateAndGet(t *testing.T) {
	ticketRepo := repository.NewTicketRepository(testDB)
	userRepo := repository.NewUserRepository(testDB)

	ctx := context.Background()

	user := &model.User{
		Name:     "Bob",
		Email:    "bob@test.com",
		Password: "hashed-password",
		Role:     model.RoleClient,
	}

	err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	// SELECT назад
	got, err := userRepo.GetUserByEmail(ctx, "bob@test.com")
	require.NoError(t, err)
	require.NotNil(t, got)

	author := got.ID
	ticket := &model.Ticket{
		Title:       "Test",
		Description: "",
		Status:      model.StatusOpen,
		Priority:    model.PriorityMedium,
		AuthorID:    author,
		AssigneeID:  &author,
	}

	err = ticketRepo.CreateTicket(ctx, ticket)
	require.NoError(t, err)
	require.NotZero(t, ticket.ID)

	res, err := ticketRepo.GetTicketByID(ctx, ticket.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, got.ID, res.AuthorID)
	assert.Equal(t, model.StatusOpen, res.Status)

}

func TestTicketRepository_Delete(t *testing.T) {
	ticketRepo := repository.NewTicketRepository(testDB)
	userRepo := repository.NewUserRepository(testDB)

	ctx := context.Background()

	user := &model.User{
		Name:     "John",
		Email:    "john@test.com",
		Password: "hashed-password",
		Role:     model.RoleClient,
	}

	err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	// SELECT назад
	got, err := userRepo.GetUserByEmail(ctx, "john@test.com")
	require.NoError(t, err)
	require.NotNil(t, got)

	author := got.ID
	ticket := &model.Ticket{
		Title:       "Test",
		Description: "leaky tap",
		Status:      model.StatusOpen,
		Priority:    model.PriorityMedium,
		AuthorID:    author,
		AssigneeID:  &author,
	}

	err = ticketRepo.CreateTicket(ctx, ticket)
	require.NoError(t, err)
	require.NotZero(t, ticket.ID)

	err = ticketRepo.DeleteTicket(ctx, ticket.ID)
	require.NoError(t, err)

	res, err := ticketRepo.GetTicketByID(ctx, ticket.ID)

	assert.ErrorIs(t, err, service.ErrTicketNotFound)
	assert.Nil(t, res)

}

func TestTicketRepository_GetByID_NotFound(t *testing.T) {
	ticketRepo := repository.NewTicketRepository(testDB)
	res, err := ticketRepo.GetTicketByID(context.Background(), 999999)
	assert.ErrorIs(t, err, service.ErrTicketNotFound)
	assert.Nil(t, res)
}
