package request

type ReqWSDiscardCards struct {
	Cards    []DiscardCards `json:"cards"`
	PlayTurn int            `json:"playTurn"`
}

type DiscardCards struct {
	CardID uint `json:"cardID"`
}
