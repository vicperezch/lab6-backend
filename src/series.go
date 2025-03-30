package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary Get list of all series.
// @Description Retrieves all series from the database.
// @Tags series
// @Produce json
// @Param search query string false "Title to search"
// @Param status query string false "Filter series by status"
// @Param sort query string false "Sort order of the series"
// @Success 200 {array} Series
// @router /series/ [get]
func getSeries(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, ranking, title, status, total_episodes, last_watched FROM series"
	sort := "asc"

	queryParams := r.URL.Query()
	if sortParam, exists := queryParams["sort"]; exists {
		if sortParam[0] == "asc" {
			sort = "desc"
		}
	}

	statusParam, statusExists := queryParams["status"]
	if statusExists && statusParam[0] != "" {
		query += " WHERE status = " + "\"" + statusParam[0] + "\""
	}

	if titleParam, exists := queryParams["search"]; exists && titleParam[0] != "" {
		if statusParam[0] != "" {
			query += " AND title = " + "\"" + titleParam[0] + "\""

		} else {
			query += " WHERE title = " + "\"" + titleParam[0] + "\""
		}
	}

	query += " ORDER BY ranking " + sort
	rows, err := db.Query(query)

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

// @Summary Get one series by ID.
// @Description Retrieves one series from the database by its ID.
// @Tags series
// @Produce json
// @Param id path int true "Series ID"
// @Success 200 {object} Series
// @router /series/{id} [get]
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

// @Summary Create a new series.
// @Description Creates a new series entry in the database.
// @Tags series
// @Accept json
// @Param series body PostRequest true "Series data"
// @Success 200 {string} Series created successfully
// @router /series/ [post]
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

// @Summary Update a series.
// @Description Changes the attributes of an existing series.
// @Tags series
// @Accept json
// @Param id path int true "Series ID"
// @Param series body Series true "Series new data"
// @Success 200 {string} Series updated successfully
// @router /series/{id} [put]
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

// @Summary Deletes a series.
// @Description Deletes an existing series entry in the database.
// @Tags series
// @Param id path int true "Series ID"
// @Success 200 {string} Series deleted successfully
// @router /series/{id} [delete]
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

// @Summary Increments an episode for a series.
// @Description Increments the last watched episode of a series by one.
// @Tags series
// @Param id path int true "Series ID"
// @Success 200 {string} Episode incremented successfully
// @router /series/{id}/episode [patch]
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

// @Summary Updates the status of one series.
// @Description Changes the status of a series to a new one.
// @Tags series
// @Param id path int true "Series ID"
// @Param series body UpdateStatusRequest true "New status"
// @Success 200 {string} Status updated successfully
// @router /series/{id}/status [patch]
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

// @Summary Upvotes one series.
// @Description Increments the rank of a series by one.
// @Tags series
// @Param id path int true "Series ID"
// @Success 200 {string} Series upvoted successfully
// @router /series/{id}/upvote [patch]
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

	if currentRanking == 1 {
		respondWithJSON(w, "Series is already rank 1")
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

// @Summary Downvotes one series.
// @Description Decrements the rank of a series by one.
// @Tags series
// @Param id path int true "Series ID"
// @Success 200 {string} Series downvoted successfully
// @router /series/{id}/downvote [patch]
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
	row := tx.QueryRow("SELECT CASE WHEN ranking = (SELECT COUNT(id) FROM series) THEN -1 ELSE ranking END AS adjusted_ranking FROM series WHERE id = ?", id)
	err = row.Scan(&currentRanking)
	if err != nil {
		respondWithError(w, "Failed to get series", http.StatusInternalServerError)
		return
	}

	if currentRanking == -1 {
		respondWithJSON(w, "Series is already at the bottom rank")
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
