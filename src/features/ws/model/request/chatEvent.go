package request

type ReqWSChat struct {
	UserID  uint   `json:"userID"`
	RoomID  uint   `json:"roomID"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
