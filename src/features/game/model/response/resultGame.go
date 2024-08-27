package response

type ResResult struct {
	Score   int      `json:"score"`
	Bonuses []string `json:"bonuses"`
}
