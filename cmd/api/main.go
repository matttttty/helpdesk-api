package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"helpdesk-api/internal/repository"
	"log"
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
}
