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
		return
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

func updateSeries(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var series Series
	if err := json.NewDecoder(r.Body).Decode(&series); err != nil {
		respondWithError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := "UPDATE series SET ranking = ?, title = ?, status = ?, total_episodes = ?, last_watched = ? WHERE id = ?"
	_, err := db.Exec(query, series.Ranking, series.Title, series.Status, series.TotalEpisodes, series.LastWatched, id)

	if err != nil {
		respondWithError(w, "Failed to update series", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Series updated successfully")
}

func deleteSeries(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		respondWithError(w, "Failed to delete series", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Series deleted successfully")
}

func incrementEpisode(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE series SET last_watched = last_watched + 1 WHERE id = ?", id)
	if err != nil {
		respondWithError(w, "Failed to increment episode", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Episode incremented successfully")
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateRequest UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		respondWithError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE series SET status = ? WHERE id = ?", updateRequest.Status, id)
	if err != nil {
		respondWithError(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Status updated successfully")
}

func upvoteSeries(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondWithError(w, "Failed to upvote", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var currentRanking int
	row := tx.QueryRow("SELECT ranking FROM series WHERE id = ?", id)
	err = row.Scan(&currentRanking)
	if err != nil {
		respondWithError(w, "Failed to get series", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = (SELECT count(id) + 1 FROM series) WHERE ranking = ? - 1", currentRanking)
	if err != nil {
		respondWithError(w, "Failed to upvote", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = ranking - 1 WHERE id = ?", id)
	if err != nil {
		respondWithError(w, "Failed to upvote", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = ? WHERE ranking = (SELECT count(id) + 1 FROM series)", currentRanking)
	if err != nil {
		respondWithError(w, "Failed to upvote", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		respondWithError(w, "Failed to write changes", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Series upvoted successfully")
}

func downvoteSeries(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		respondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondWithError(w, "Failed to downvote", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var currentRanking int
	row := tx.QueryRow("SELECT ranking FROM series WHERE id = ?", id)
	err = row.Scan(&currentRanking)
	if err != nil {
		respondWithError(w, "Failed to get series", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = (SELECT count(id) + 1 FROM series) WHERE ranking = ? + 1", currentRanking)
	if err != nil {
		respondWithError(w, "Failed to downvote", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = ranking + 1 WHERE id = ?", id)
	if err != nil {
		respondWithError(w, "Failed to downvote", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE series SET ranking = ? WHERE ranking = (SELECT count(id) + 1 FROM series)", currentRanking)
	if err != nil {
		respondWithError(w, "Failed to downvote", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		respondWithError(w, "Failed to write changes", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "Series upvoted successfully")
}
