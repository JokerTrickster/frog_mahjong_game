package entity

type ResultEntitySQL struct {
	UserID uint
	RoomID uint
	Cards  []uint
}

type ResultEntity struct {
	Cards []ResultCard `json:"cards"`
}

type ResultCard struct {
	CardID uint   `json:"cardID"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type ResultCardEntity struct {
	CardID uint   `json:"cardID"`
	Count  int    `json:"count"`
	Color  string `json:"color"`
	State  string `json:"state"`
}
