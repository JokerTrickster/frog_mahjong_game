package request

type ReqWSItemChange struct {
	RoomID uint `json:"roomID omitempty"`
	UserID uint `json:"userID"`
	ItemID uint `json:"itemID"`
}
