package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"helpdesk-api/internal/handler"
	"helpdesk-api/internal/repository"
	"helpdesk-api/internal/service"
	"log"
	"net/http"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Helpdesk API")
	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("Connected to database", db.Stats())

	h := handler.NewAuthHandler(service.NewUserService(repository.NewUserRepository(db)))
	router := handler.NewRouter(h)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}
