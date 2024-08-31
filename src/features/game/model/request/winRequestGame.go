package request

type ReqWinRequest struct {
	UserID uint `json:"userID"`
	RoomID uint `json:"roomID"`
	Score  int  `json:"score"`
}
