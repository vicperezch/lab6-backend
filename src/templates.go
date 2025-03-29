package main

type Series struct {
	Id	int `json:"id"`
	Ranking int `json:"ranking"`
	Title	string `json:"title"`
	Status	string `json:"status"`
	TotalEpisodes	int `json:"totalEpisodes"`
	LastWatched	int `json:"lastEpisodeWatched"`
}

type PostRequest struct {
	Title string `json:"title"`
	Status string `json:"status"`
	LastWatched int `json:"lastEpisodeWatched"`
	TotalEpisodes int `json:"totalEpisodes"`
	Ranking int `json:"ranking"`
}
