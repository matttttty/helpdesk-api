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

func isStaff(role model.Role) bool {
	return role == model.RoleAdmin || role == model.RoleAgent
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

func (s *TicketService) GetTicketByID(ctx context.Context, id int64, userID int64, role model.Role) (*model.Ticket, error) {

	t, err := s.repo.GetTicketByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isStaff(role) && userID != t.AuthorID {
		return nil, ErrTicketNotFound

	}

	return t, nil

}

func (s *TicketService) GetAllTickets(ctx context.Context, userID int64, role model.Role) ([]*model.Ticket, error) {

	if !isStaff(role) {
		return s.repo.GetTicketsByAuthorID(ctx, userID)
	}

	return s.repo.GetAllTickets(ctx)

}

func (s *TicketService) UpdateTicket(ctx context.Context, ticket *model.Ticket, userID int64, role model.Role) error {
	existing, err := s.repo.GetTicketByID(ctx, ticket.ID)
	if err != nil {
		return err
	}

	if !isStaff(role) && existing.AuthorID != userID {
		return ErrTicketNotFound
	}

	if ticket.Title == "" {
		return ErrTitleRequired
	}

	ticket.AuthorID = existing.AuthorID
	return s.repo.UpdateTicket(ctx, ticket)

}

func (s *TicketService) DeleteTicket(ctx context.Context, id int64, role model.Role) error {
	if role != model.RoleAdmin {
		return ErrAccessDenied
	}
	return s.repo.DeleteTicket(ctx, id)

}
