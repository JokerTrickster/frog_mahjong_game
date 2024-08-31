package entity

type WSLoanEntity struct {
	RoomID       uint `json:"roomID omitempty"`
	UserID       uint `json:"userID"`
	CardID       uint `json:"cardID"`
	TargetUserID uint `json:"targetUserID"`
}
