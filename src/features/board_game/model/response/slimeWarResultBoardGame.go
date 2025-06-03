package response

type ResSlimeWarResult struct {
	Users []SlimeWarResult `json:"users"`
}

type SlimeWarResult struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
	Result int `json:"result"`
}
