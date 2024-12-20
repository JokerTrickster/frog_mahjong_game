package entity

type WSAbnormalEntity struct {
	RoomID         uint `json:"roomID omitempty"`
	UserID         uint `json:"userID omitempty"`
	AbnormalUserID uint `json:"abnormalUserID"`
}
