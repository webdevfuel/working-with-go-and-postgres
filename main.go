package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := OpenDatabase()
	if err != nil {
		log.Printf("error opening database connection %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		err := DB.QueryRow("INSERT INTO books (name, description) VALUES ('The Greatest Book Ever', '');").Err()
		if err != nil {
			log.Printf("error insert book into books table %v", err)
		}
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe("localhost:3000", r)
}
