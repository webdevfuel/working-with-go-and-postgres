package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	rows, err := DB.Query("SELECT id, name, description FROM books;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error querying books table %v", err)
		return
	}

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Description); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("scanning books table rows into struct Book %v", err)
			return
		}
		books = append(books, book)
	}

	j, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling books into json", err)
		return
	}

	w.Write(j)
}
