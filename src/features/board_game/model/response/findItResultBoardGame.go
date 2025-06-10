package response

type ResFindItResult struct {
	Users []FindItResult `json:"users"`
}

type FindItResult struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
	Result int `json:"result"`
}
