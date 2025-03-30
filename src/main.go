package main

import (
	"database/sql"
	"log"
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
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Get("/api/series", getSeries)
	router.Get("/api/series/{id}", getSeriesById)
	router.Post("/api/series", createSeries)
	router.Put("/api/series/{id}", updateSeries)
	router.Delete("/api/series/{id}", deleteSeries)
	router.Patch("/api/series/{id}/episode", incrementEpisode)
	router.Patch("/api/series/{id}/status", updateStatus)
	router.Patch("/api/series/{id}/upvote", upvoteSeries)
	router.Patch("/api/series/{id}/downvote", downvoteSeries)

	http.ListenAndServe(":8080", router)
}
