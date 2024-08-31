package entity

type WSChatEntity struct {
	UserID  uint   `json:"userID"`
	RoomID  uint   `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
