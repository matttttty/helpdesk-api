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
	return m.Called(ctx, ticket).Error(0)
}
func (m *MockTicketRepo) GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error) {
	args := m.Called(ctx, id)
	t, _ := args.Get(0).(*model.Ticket)
	return t, args.Error(1)
}
func (m *MockTicketRepo) GetAllTickets(ctx context.Context) ([]*model.Ticket, error) {
	args := m.Called(ctx)
	t, _ := args.Get(0).([]*model.Ticket)
	return t, args.Error(1)
}
func (m *MockTicketRepo) GetTicketsByAuthorID(ctx context.Context, authorID int64) ([]*model.Ticket, error) {
	args := m.Called(ctx, authorID)
	t, _ := args.Get(0).([]*model.Ticket)
	return t, args.Error(1)
}
func (m *MockTicketRepo) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {
	return m.Called(ctx, ticket).Error(0)
}
func (m *MockTicketRepo) DeleteTicket(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreateTicket_EmptyTitle(t *testing.T) {
	repo := new(MockTicketRepo)
	svc := service.NewTicketService(repo)

	err := svc.CreateTicket(context.Background(), &model.Ticket{Title: ""})

	assert.ErrorIs(t, err, service.ErrTitleRequired)
	repo.AssertNotCalled(t, "CreateTicket") // у repo навіть не дійшло
}

func TestCreateTicket_SetsDefaults(t *testing.T) {
	repo := new(MockTicketRepo)
	// ловимо тикет, який service передав у repo, і перевіряємо дефолти
	repo.On("CreateTicket", mock.Anything, mock.MatchedBy(func(tk *model.Ticket) bool {
		return tk.Status == model.StatusOpen && tk.Priority == model.PriorityMedium
	})).Return(nil)
	svc := service.NewTicketService(repo)

	err := svc.CreateTicket(context.Background(), &model.Ticket{Title: "test"})

	assert.NoError(t, err)
	repo.AssertExpectations(t) // впаде, якщо дефолти не виставились
}

func TestDeleteTicket_AccessDenied(t *testing.T) {
	repo := new(MockTicketRepo)
	svc := service.NewTicketService(repo)

	err := svc.DeleteTicket(context.Background(), 1, model.RoleClient)

	assert.ErrorIs(t, err, service.ErrAccessDenied)
	repo.AssertNotCalled(t, "DeleteTicket")
}

func TestDeleteTicket_AdminSucceeds(t *testing.T) {
	repo := new(MockTicketRepo)
	repo.On("DeleteTicket", mock.Anything, int64(1)).Return(nil)
	svc := service.NewTicketService(repo)

	err := svc.DeleteTicket(context.Background(), 1, model.RoleAdmin)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetTicketByID_Access(t *testing.T) {
	const ownerID, strangerID, adminID, agentID = 1, 2, 99, 50

	cases := []struct {
		name    string
		userID  int64
		role    model.Role
		wantErr error
	}{
		{"owner reads own", ownerID, model.RoleClient, nil},
		{"admin reads foreign", adminID, model.RoleAdmin, nil},
		{"agent reads foreign", agentID, model.RoleAgent, nil},
		{"stranger reads foreign", strangerID, model.RoleClient, service.ErrTicketNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockTicketRepo)
			repo.On("GetTicketByID", mock.Anything, int64(10)).
				Return(&model.Ticket{ID: 10, AuthorID: ownerID, Title: "x"}, nil)
			svc := service.NewTicketService(repo)

			_, err := svc.GetTicketByID(context.Background(), 10, tc.userID, tc.role)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUpdateTicket_Access(t *testing.T) {
	const ownerID, strangerID, adminID = 1, 2, 99

	cases := []struct {
		name    string
		userID  int64
		role    model.Role
		wantErr error
	}{
		{"admin edits any", adminID, model.RoleAdmin, nil},
		{"owner edits own", ownerID, model.RoleClient, nil},
		{"stranger edits foreign", strangerID, model.RoleClient, service.ErrTicketNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockTicketRepo)
			repo.On("GetTicketByID", mock.Anything, int64(10)).
				Return(&model.Ticket{ID: 10, AuthorID: ownerID, Title: "old"}, nil)
			// UpdateTicket у repo кличеться лише коли доступ є
			repo.On("UpdateTicket", mock.Anything, mock.Anything).Return(nil).Maybe()
			svc := service.NewTicketService(repo)

			err := svc.UpdateTicket(context.Background(),
				&model.Ticket{ID: 10, Title: "new"}, tc.userID, tc.role)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUpdateTicket_PreservesAuthorID(t *testing.T) {
	const ownerID int64 = 1
	repo := new(MockTicketRepo)
	repo.On("GetTicketByID", mock.Anything, int64(10)).
		Return(&model.Ticket{ID: 10, AuthorID: ownerID, Title: "old"}, nil)
	// ключове: тикет, що йде в repo, мусить зберегти AuthorID, а не занулити
	repo.On("UpdateTicket", mock.Anything, mock.MatchedBy(func(tk *model.Ticket) bool {
		return tk.AuthorID == ownerID
	})).Return(nil)
	svc := service.NewTicketService(repo)

	// навмисно НЕ передаємо AuthorID у вхідному тикеті — service має підтягнути зі старого
	err := svc.UpdateTicket(context.Background(),
		&model.Ticket{ID: 10, Title: "new"}, ownerID, model.RoleClient)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetAllTickets_StaffSeesAll(t *testing.T) {
	repo := new(MockTicketRepo)
	repo.On("GetAllTickets", mock.Anything).Return([]*model.Ticket{}, nil)
	svc := service.NewTicketService(repo)

	_, err := svc.GetAllTickets(context.Background(), 50, model.RoleAgent)

	assert.NoError(t, err)
	repo.AssertCalled(t, "GetAllTickets", mock.Anything)
	repo.AssertNotCalled(t, "GetTicketsByAuthorID")
}

func TestGetAllTickets_ClientSeesOwn(t *testing.T) {
	const clientID int64 = 7
	repo := new(MockTicketRepo)
	repo.On("GetTicketsByAuthorID", mock.Anything, clientID).Return([]*model.Ticket{}, nil)
	svc := service.NewTicketService(repo)

	_, err := svc.GetAllTickets(context.Background(), clientID, model.RoleClient)

	assert.NoError(t, err)
	repo.AssertCalled(t, "GetTicketsByAuthorID", mock.Anything, clientID)
	repo.AssertNotCalled(t, "GetAllTickets")
}
