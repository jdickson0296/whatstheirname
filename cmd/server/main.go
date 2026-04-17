package main

import (
	"log"
	"net/http"

	"github.com/jdickson0296/whatstheirname/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	http.HandleFunc("/search", handlers.SearchHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
