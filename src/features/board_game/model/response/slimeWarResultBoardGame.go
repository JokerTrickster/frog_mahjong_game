package response

type ResSlimeWarResult struct {
	Result []SlimeWarResult `json:"result"`
}

type SlimeWarResult struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
}
