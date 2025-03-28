package main

import (
	"log"
	"net/http"
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
