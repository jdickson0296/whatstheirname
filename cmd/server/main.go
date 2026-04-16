package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Search Works!")
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}