package service

import (
	"context"
	"errors"
	"helpdesk-api/internal/model"
	"helpdesk-api/internal/repository"
)

type ticketService struct {
	repo *repository.TicketRepository
}

func NewTicketService(repo *repository.TicketRepository) *ticketService {
	return &ticketService{repo: repo}
}

func (s *ticketService) CreateTicket(ctx context.Context, ticket *model.Ticket) error {
	if ticket.Title == "" {
		return errors.New("title is required")
	}
	if ticket.Status == "" {
		ticket.Status = model.StatusOpen
	}
	if ticket.Priority == "" {
		ticket.Priority = model.PriorityMedium
	}
	return s.repo.CreateTicket(ctx, ticket)
}

func (s *ticketService) GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error) {

	return s.repo.GetTicketByID(ctx, id)

}

func (s *ticketService) GetAllTickets(ctx context.Context) ([]*model.Ticket, error) {

	return s.repo.GetAllTickets(ctx)

}

func (s *ticketService) GetTicketsByAuthorID(ctx context.Context, AuthorID int64) ([]*model.Ticket, error) {

	return s.repo.GetTicketsByAuthorID(ctx, AuthorID)

}

func (s *ticketService) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {

	return s.repo.UpdateTicket(ctx, ticket)

}

func (s *ticketService) DeleteTicket(ctx context.Context, id int64) error {

	return s.repo.DeleteTicket(ctx, id)

}
