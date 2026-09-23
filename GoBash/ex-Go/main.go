package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hello, Go server!")
	})

	log.Println("server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
