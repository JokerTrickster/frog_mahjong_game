package request

// card 요청 구조체 만들어줘

type ReqOwnership struct {
	Cards []Card `json:"cards"`
}

/*
roomID int
userID int
cardID int
state string (owned 카드 상태)
*/
type Card struct {
	CardID int    `json:"cardID"`
	State  string `json:"state"`
	RoomID int    `json:"roomID"`
	UserID int    `json:"userID"`
}
