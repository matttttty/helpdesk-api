package service

import (
	"context"
	"helpdesk-api/internal/model"
)

type TicketRepo interface {
	CreateTicket(ctx context.Context, ticket *model.Ticket) error
	GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error)
	GetAllTickets(ctx context.Context) ([]*model.Ticket, error)
	GetTicketsByAuthorID(ctx context.Context, AuthorID int64) ([]*model.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *model.Ticket) error
	DeleteTicket(ctx context.Context, id int64) error
}

type TicketService struct {
	repo TicketRepo
}

func NewTicketService(repo TicketRepo) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) CreateTicket(ctx context.Context, ticket *model.Ticket) error {
	if ticket.Title == "" {
		return ErrTitleRequired
	}
	if ticket.Status == "" {
		ticket.Status = model.StatusOpen
	}
	if ticket.Priority == "" {
		ticket.Priority = model.PriorityMedium
	}
	return s.repo.CreateTicket(ctx, ticket)
}

func (s *TicketService) GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error) {

	return s.repo.GetTicketByID(ctx, id)

}

func (s *TicketService) GetAllTickets(ctx context.Context) ([]*model.Ticket, error) {

	return s.repo.GetAllTickets(ctx)

}

func (s *TicketService) GetTicketsByAuthorID(ctx context.Context, AuthorID int64) ([]*model.Ticket, error) {

	return s.repo.GetTicketsByAuthorID(ctx, AuthorID)

}

func (s *TicketService) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {

	return s.repo.UpdateTicket(ctx, ticket)

}

func (s *TicketService) DeleteTicket(ctx context.Context, id int64, role model.Role) error {
	if role != model.RoleAdmin {
		return ErrAccessDenied
	}
	return s.repo.DeleteTicket(ctx, id)

}
