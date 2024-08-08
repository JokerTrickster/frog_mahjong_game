package request

type ReqWSDora struct {
	Cards    []Card `json:"cards"`
	PlayTurn int    `json:"playTurn"`
}

type Card struct {
	CardID uint `json:"cardID"`
}
