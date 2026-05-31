package service_test

import (
	"context"
	"testing"

	"helpdesk-api/internal/model"
	"helpdesk-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTicketRepo struct {
	mock.Mock
}

func (m *MockTicketRepo) CreateTicket(ctx context.Context, ticket *model.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepo) GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.Ticket)
	return u, args.Error(1)
}

func (m *MockTicketRepo) GetAllTickets(ctx context.Context) ([]*model.Ticket, error) {
	args := m.Called(ctx)
	u, _ := args.Get(0).([]*model.Ticket)
	return u, args.Error(1)
}

func (m *MockTicketRepo) GetTicketsByAuthorID(ctx context.Context, AuthorID int64) ([]*model.Ticket, error) {
	args := m.Called(ctx, AuthorID)
	u, _ := args.Get(0).([]*model.Ticket)
	return u, args.Error(1)
}

func (m *MockTicketRepo) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepo) DeleteTicket(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTicket_EmptyTitle(t *testing.T) {
	repo := new(MockTicketRepo)

	svc := service.NewTicketService(repo)

	err := svc.CreateTicket(context.Background(), &model.Ticket{Title: ""})
	assert.ErrorIs(t, err, service.ErrTitleRequired)
	repo.AssertExpectations(t)
}

func TestCreateTicket_SetsDefaults(t *testing.T) {
	repo := new(MockTicketRepo)
	repo.On("CreateTicket", mock.Anything, mock.Anything).Return(nil)

	svc := service.NewTicketService(repo)

	err := svc.CreateTicket(context.Background(), &model.Ticket{Title: "test", Status: "", Priority: ""})
	assert.NoError(t, err)
	assert.Equal(t, model.StatusOpen, model.StatusOpen)
	assert.Equal(t, model.PriorityMedium, model.PriorityMedium)
	repo.AssertExpectations(t)
}

func TestDeleteTicket_AccessDenied(t *testing.T) {
	repo := new(MockTicketRepo)

	svc := service.NewTicketService(repo)

	err := svc.DeleteTicket(context.Background(), 1, model.RoleClient)
	assert.ErrorIs(t, err, service.ErrAccessDenied)
	repo.AssertExpectations(t)
}

func TestDeleteTicket_AdminSucceeds(t *testing.T) {
	repo := new(MockTicketRepo)

	repo.On("DeleteTicket", mock.Anything, int64(1)).Return(nil)

	svc := service.NewTicketService(repo)

	err := svc.DeleteTicket(context.Background(), 1, model.RoleAdmin)

	assert.NoError(t, err)
	repo.AssertExpectations(t)

}
