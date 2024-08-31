package response

type ResScoreCalculate struct {
	Score   int      `json:"score"`
	Bonuses []string `json:"bonuses"`
}
