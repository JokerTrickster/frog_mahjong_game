package response

type ResFindItRankBoardGame struct {
	RankUserList []RankUser `json:"rankUserList"`
}

type RankUser struct {
	UserID    int    `json:"userID"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
	Rank      int    `json:"rank"`
	ProfileID int    `json:"profileID"`
}
