package main

type Series struct {
	Id	int `json:"id"`
	Ranking int `json:"ranking"`
	Title	string `json:"title"`
	Status	string `json:"status"`
	TotalEpisodes	int `json:"totalEpisodes"`
	LastWatched	int `json:"lastEpisodeWatched"`
}
