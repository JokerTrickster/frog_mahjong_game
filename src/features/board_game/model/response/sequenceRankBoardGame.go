package response

type ResSequenceRank struct {
	RankUserList []RankUser `json:"rankUserList"`
}

type SequenceRank struct {
	UserID int `json:"userID"`
	Score  int `json:"score"`
}
