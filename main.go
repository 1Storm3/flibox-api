package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hello, welcome to the Kinopoisk API!")
	})

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
