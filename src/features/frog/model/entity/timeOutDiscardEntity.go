package entity

type WSTimeOutDiscardCardsEntity struct {
	RoomID uint `json:"roomID omitempty"`
	UserID uint `json:"userID"`
	CardID uint `json:"cardID"`
}
