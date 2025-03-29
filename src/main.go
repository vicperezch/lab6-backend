package main

import (
	"log"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "../db/db.sqlite")

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	router.Get("/api/series", getSeries)
	router.Post("/api/series", createSeries)

	http.ListenAndServe(":8080", router)
}
