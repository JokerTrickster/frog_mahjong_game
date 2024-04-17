package response

type ResUserListRoom struct {
	Users []User `json:"users"`
}
type User struct {
	UserID         int    `json:"userID" gorm:"column:user_id"`
	RoomUserID     int    `json:"roomUserID" grom:"column:room_user_id"`
	PlayerState    string `json:"playerState" gorm:"column:player_state"`
	TurnNumber     int    `json:"turnNumber" gorm:"column:turn_number"`
	OwnedCardCount int    `json:"ownedCardCount" gorm:"column:owned_card_count"`
	RoomID         int    `json:"roomID" gorm:"column:room_id"`
	Score          int    `json:"score" gorm:"column:score"`
	UserName       string `json:"userName" gorm:"column:user_name"`
	UserEmail      string `json:"userEmail" gorm:"column:user_email"`
	Owner          bool   `json:"owner"`
}
