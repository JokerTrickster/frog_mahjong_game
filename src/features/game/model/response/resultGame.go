package response

type ResResult struct {
	Score   int      `json:"score"`
	Winner  uint64   `json:"winner"`
	Bonuses []string `json:"bonuses"`
}
