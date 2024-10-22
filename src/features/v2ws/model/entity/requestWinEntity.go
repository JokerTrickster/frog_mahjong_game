package entity

type WSRequestWinEntity struct {
	Cards []RequestWinCard `json:"cards"`
}

type RequestWinCard struct {
	CardID uint   `json:"cardID"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
