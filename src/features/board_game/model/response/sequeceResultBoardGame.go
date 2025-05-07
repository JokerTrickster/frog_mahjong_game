package response

type ResSequenceResult struct {
	Result []SequenceResult `json:"result"`
}

type SequenceResult struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
}
