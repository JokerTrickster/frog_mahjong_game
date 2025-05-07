package request

type ReqGameOverBoardGame struct {
	Winner   bool `json:"winner"`
	GameType int  `json:"gameType"`
	RoomID   int  `json:"roomID"`
}
