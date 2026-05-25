package handler

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"helpdesk-api/pkg/middleware"
)

func NewRouter(authHandler *AuthHandler, ticketHandler *TicketHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Route("/tickets", func(r chi.Router) {
			r.Get("/", ticketHandler.GetAllTickets)
			r.Post("/", ticketHandler.CreateTicket)
			r.Get("/{id}", ticketHandler.GetTicketByID)
			r.Put("/{id}", ticketHandler.UpdateTicket)
			r.Delete("/{id}", ticketHandler.DeleteTicket)
		})
	})

	return r
}
