package request

type ReqNextTurn struct {
	RoomID      int    `json:"roomID"`
	UserID      int    `json:"userID"`
	PlayerState string `json:"playerState"`
	TurnNumber  int    `json:"turnNumber"`
}
