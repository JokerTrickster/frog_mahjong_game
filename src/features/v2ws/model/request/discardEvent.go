package request

type ReqWSDiscardCards struct {
	CardID   uint `json:"cardID"`
	PlayTurn int  `json:"playTurn"`
}
