package main

import (
	"fmt"
	"kinopoisk-api/internal/storage/postgres"
	"log"
	"net/http"
)

func main() {

	connStr := ""

	storage, err := postgres.NewStorage(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(storage *postgres.Storage) {
		err = storage.Close()
		if err != nil {

		}
	}(storage)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err = fmt.Fprintln(w, "Hello, welcome to the Kinopoisk API!")
		if err != nil {
			return
		}
		_, err = w.Write([]byte("Hello, welcome to the Kinopoisk API!"))
		if err != nil {
			return
		}
	})

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
