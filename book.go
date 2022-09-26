package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
