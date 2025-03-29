package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getSeries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, ranking, title, status, total_episodes, last_watched FROM series")
	
	if err != nil {
		log.Fatal("Failed to get all series:", err)
		return
	}

	defer rows.Close()

	var series []Series
	for rows.Next() {
		var s Series
		rows.Scan(&s.Id, &s.Ranking, &s.Title, &s.Status, &s.TotalEpisodes, &s.LastWatched)
		series = append(series, s)
	}

	respondWithJSON(w, series)
}

func getSeriesById(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	
	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	row := db.QueryRow("SELECT id, ranking, title, status, total_episodes, last_watched FROM series WHERE id = ?", id)

	var s Series
	err := row.Scan(&s.Id, &s.Ranking, &s.Title, &s.Status, &s.TotalEpisodes, &s.LastWatched)
	if err != nil {
		respondWithError(w, "Failed to get series", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, s)
}

func createSeries(w http.ResponseWriter, r *http.Request) {
	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request", http.StatusBadRequest)		
	}

	_, err := db.Exec(
		"INSERT INTO series (ranking, title, status, total_episodes, last_watched) VALUES (?, ?, ?, ?, ?)",
		req.Ranking, req.Title, req.Status, req.TotalEpisodes, req.LastWatched,
	)

	if err != nil {
		log.Println("Error creating series: ", err)
		respondWithError(w, "Error creating series.", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Series created successfully")
}
