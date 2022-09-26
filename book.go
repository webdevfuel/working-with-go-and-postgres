package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateBookBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func create(w http.ResponseWriter, r *http.Request) {
	var body CreateBookBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body into CreateBookBody struct %v", err)
		return
	}

	if err := DB.QueryRow("INSERT INTO books (name, description) VALUES ($1, $2)", body.Name, body.Description).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error inserting book into books table %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getAll(w http.ResponseWriter, _ *http.Request) {
	var books []Book
	if err := DB.Select(&books, "SELECT id, name, description FROM books;"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error querying books table %v", err)
		return
	}

	j, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling books into json %v", err)
		return
	}

	w.Write(j)
}

func get(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error parsing %d from string into integer %v", bookID, err)
		return
	}

	var book Book
	if err := DB.QueryRow("SELECT id, name, description FROM books WHERE id = $1;", bookID).Scan(&book.ID, &book.Name, &book.Description); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error querying book from books table with id %d %v", bookID, err)
		return
	}

	j, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling books into json %v", err)
		return
	}

	w.Write(j)
}
