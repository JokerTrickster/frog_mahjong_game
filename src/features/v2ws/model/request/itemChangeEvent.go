package request

type ReqWSItemChange struct {
	RoomID uint `json:"roomID omitempty"`
	UserID uint `json:"userID"`
	CardID uint `json:"cardID"`
}
