package entity

type WSImportSingleCardEntity struct {
	RoomID uint `json:"roomID omitempty"`
	UserID uint `json:"userID"`
	CardID uint `json:"cardID"`
}
