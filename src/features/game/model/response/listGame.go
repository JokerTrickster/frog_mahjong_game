package response

type ResListGame struct {
	Games      []GameInfo `json:"games"`
	TotalCount int        `json:"totalCount"`
}

type GameInfo struct {
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	HashTag     string `json:"hashTag"`
	Category    string `json:"category"`
	IsEnabled   bool   `json:"isEnabled"`
	YoutubeUrl  string `json:"youtubeUrl"`
}
