package entity

type V2WSRequestWinEntity struct {
	RoomID uint  `json:"roomID omitempty"`
	UserID uint  `json:"userID"`
	Cards  []int `json:"cards"`
}
