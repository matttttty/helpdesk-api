package repository

import (
	"context"
	"database/sql"
	"fmt"
	"helpdesk-api/internal/model"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTicket(ctx context.Context, ticket *model.Ticket) error {

	query := `INSERT INTO tickets (title, description, status, priority, author_id, assignee_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, ticket.Title, ticket.Description, ticket.Status, ticket.Priority, ticket.AuthorID, ticket.AssigneeID).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TicketRepository) GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error) {

	query := `SELECT id, title, description, status, priority, author_id, assignee_id, created_at, updated_at FROM tickets WHERE id = $1`

	ticket := &model.Ticket{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.Priority, &ticket.AuthorID, &ticket.AssigneeID, &ticket.CreatedAt, &ticket.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("GetTicketById: %w", err)
	}

	return ticket, nil
}

func (r *TicketRepository) GetAllTickets(ctx context.Context) ([]*model.Ticket, error) {

	query := `SELECT id, title, description, status, priority, author_id, assignee_id, created_at, updated_at FROM tickets`

	tikests := []*model.Ticket{}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllTickets: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		ticket := &model.Ticket{}
		if err := rows.Scan(
			&ticket.ID, &ticket.Title, &ticket.Description,
			&ticket.Status, &ticket.Priority, &ticket.AuthorID,
			&ticket.AssigneeID, &ticket.CreatedAt, &ticket.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetAllTickets: %w", err)
		}
		tikests = append(tikests, ticket)
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, fmt.Errorf("GetAllTickets: %w", rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllTickets: %w", err)
	}
	return tikests, nil
}

func (r *TicketRepository) GetTicketsByAuthorID(ctx context.Context, AuthorID int64) ([]*model.Ticket, error) {

	query := `SELECT id, title, description, status, priority, author_id, assignee_id, created_at, updated_at FROM tickets WHERE author_id = $1`

	tickets := []*model.Ticket{}

	rows, err := r.db.QueryContext(ctx, query, AuthorID)
	if err != nil {
		return nil, fmt.Errorf("GetTicketsByAuthorID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		ticket := model.Ticket{}
		if err := rows.Scan(
			&ticket.ID, &ticket.Title, &ticket.Description,
			&ticket.Status, &ticket.Priority, &ticket.AuthorID,
			&ticket.AssigneeID, &ticket.CreatedAt, &ticket.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetTicketsByAuthorID: %w", err)
		}
		tickets = append(tickets, &ticket)
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, fmt.Errorf("GetTicketsByAuthorID: %w", rerr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTicketsByAuthorID: %w", err)
	}
	return tickets, nil
}

func (r *TicketRepository) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {

	query := `UPDATE tickets SET title=$1, description=$2, status=$3, priority=$4, assignee_id=$5, updated_at=NOW() WHERE id=$6 RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		ticket.Title, ticket.Description, ticket.Status,
		ticket.Priority, ticket.AssigneeID, ticket.ID,
	).Scan(&ticket.UpdatedAt)
	if err != nil {
		return fmt.Errorf("UpdateTicket: %w", err)
	}

	return nil
}

func (r *TicketRepository) DeleteTicket(ctx context.Context, id int64) error {

	query := `DELETE FROM tickets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteTicket: %w", err)
	}
	return nil
}
