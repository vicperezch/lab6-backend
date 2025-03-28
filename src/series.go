package main

import (
	"net/http"
)

func getSeries(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, ranking, title, status, total_episodes, last_watched FROM series")
	defer rows.Close()

	var series []Series

	for rows.Next() {
		var s Series
		rows.Scan(&s.Id, &s.Ranking, &s.Title, &s.Status, &s.TotalEpisodes, &s.LastWatched)
		series = append(series, s)
	}

	respondWithJSON(w, series)
}
