package request

type ReqDora struct {
	RoomID    int    `json:"roomID"`
	CardName  string `json:"cardName"`
	CardColor string `json:"cardColor"`
	CardState string `json:"cardState"`
}
