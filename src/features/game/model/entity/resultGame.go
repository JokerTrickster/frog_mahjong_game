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
