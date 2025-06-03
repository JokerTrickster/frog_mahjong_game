package request

type ReqGameOverBoardGame struct {
	GameType int `json:"gameType"`
	RoomID   int `json:"roomID"`
	UserID   int `json:"userID"`
	Score    int `json:"score"`
	Result   int `json:"result"`
}
