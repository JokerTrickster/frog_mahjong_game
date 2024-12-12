package entity

type WSItemChangeEntity struct {
	RoomID uint `json:"roomID omitempty"`
	UserID uint `json:"userID"`
	ItemID uint `json:"itemID"`
}
