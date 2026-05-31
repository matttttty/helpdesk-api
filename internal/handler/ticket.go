package handler

import (
	"encoding/json"
	"helpdesk-api/internal/model"
	"helpdesk-api/internal/service"
	"helpdesk-api/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: service}
}

func (s *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var ticket model.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ticket.Title == "" {
		http.Error(w, "Title should be filled", http.StatusBadRequest)
		return
	}
	if ticket.Description == "" {
		http.Error(w, "Description should be filled", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}
	ticket.AuthorID = claims.UserID

	err = s.ticketService.CreateTicket(r.Context(), &ticket)
	if err != nil {
		writeError(w, http.StatusBadRequest, "recheck the fields content", nil)
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	ticket, err := h.ticketService.GetAllTickets(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch the tickets", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ticket, err := h.ticketService.GetTicketByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch the ticket", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var ticket *model.Ticket
	err = json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ticket.ID = id
	err = h.ticketService.UpdateTicket(r.Context(), ticket)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update ticket", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err = h.ticketService.DeleteTicket(r.Context(), id, model.Role(claims.Role))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete ticket", err)
		return
	}

	writeJSON(w, http.StatusOK, "Ticket deleted")
}
