package main

import (
	"context"
	"errors"
	"helpdesk-api/internal/config"
	"helpdesk-api/internal/handler"
	"helpdesk-api/internal/repository"
	"helpdesk-api/internal/service"
	"helpdesk-api/pkg/middleware"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// @title           Helpdesk API
// @version         1.0
// @description     Ticket management API with JWT auth
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	middleware.LoadJWTSecret(cfg.JWTSecret)
	log.Println("Helpdesk API starting")
	db, err := repository.NewDB(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	log.Println("connected to database")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	ticketRepo := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	router := handler.NewRouter(authHandler, ticketHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}

	}()

	select {
	case err := <-serverErr:
		log.Fatalf("server failed: %v", err)
	case <-ctx.Done():
		log.Println("shutdown signal received, draining requests...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("forced shutdown: %v", err)
	} else {
		log.Println("server stopped gracefully")
	}

}
