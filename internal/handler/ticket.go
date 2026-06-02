package handler

import (
	"encoding/json"
	"errors"
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

// CreateTicket godoc
// @Summary      Create a ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticket  body      model.Ticket  true  "Ticket payload"
// @Success      201     {object}  model.Ticket
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /tickets [post]
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket model.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if ticket.Description == "" {
		writeError(w, http.StatusBadRequest, "Description should be filled", nil)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "User is not authenticated", nil)
		return
	}
	ticket.AuthorID = claims.UserID

	err = h.ticketService.CreateTicket(r.Context(), &ticket)
	if err != nil {
		if errors.Is(err, service.ErrTitleRequired) {
			writeError(w, http.StatusBadRequest, "title is required", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create ticket", err)
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

// GetAllTickets godoc
// @Summary      Get all tickets
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200     {array}   model.Ticket
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /tickets [get]
func (h *TicketHandler) GetAllTickets(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	ticket, err := h.ticketService.GetAllTickets(r.Context(), claims.UserID, model.Role(claims.Role))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch the tickets", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

// GetTicketByID godoc
// @Summary      Get a ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param  id  path  int  true  "Ticket ID"
// @Success      200     {object}  model.Ticket
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /tickets/{id} [get]
func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	ticket, err := h.ticketService.GetTicketByID(r.Context(), id, claims.UserID, model.Role(claims.Role))
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			writeError(w, http.StatusNotFound, "ticket not found", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch the ticket", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

// UpdateTicket godoc
// @Summary      Update a ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticket  body      model.Ticket  true "Ticket payload"
// @Success      200     {object}  model.Ticket
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /tickets/{id} [put]
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	var ticket model.Ticket
	err = json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	ticket.ID = id
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	err = h.ticketService.UpdateTicket(r.Context(), &ticket, claims.UserID, model.Role(claims.Role))
	if err != nil {
		if errors.Is(err, service.ErrTitleRequired) {
			writeError(w, http.StatusBadRequest, "title is required", nil)
			return
		}
		if errors.Is(err, service.ErrTicketNotFound) {
			writeError(w, http.StatusNotFound, "ticket not found", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update a ticket", err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

// DeleteTicket godoc
// @Summary      Delete a ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param  id  path  int  true  "Ticket ID"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /tickets/{id} [delete]
func (h *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	err = h.ticketService.DeleteTicket(r.Context(), id, model.Role(claims.Role))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccessDenied):
			writeError(w, http.StatusForbidden, "access denied", nil) // 403
		case errors.Is(err, service.ErrTicketNotFound):
			writeError(w, http.StatusNotFound, "ticket not found", nil) // 404
		default:
			writeError(w, http.StatusInternalServerError, "failed to delete ticket", err)
		}
		return
	}

	writeJSON(w, http.StatusOK, "Ticket deleted")
}
