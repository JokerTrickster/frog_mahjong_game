package request

type ReqDiscard struct {
	CardID    int    `json:"cardID"`
	UserID    int    `json:"userID"`
	RoomID    int    `json:"roomID"`
	CardState string `json:"cardState"`
}
