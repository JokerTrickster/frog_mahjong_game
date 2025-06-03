package response

type ResSequenceResult struct {
	Users []SequenceResult `json:"users"`
}

type SequenceResult struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
	Result int `json:"result"`
}
