package request

type ReqWSImportCards struct {
	Cards    []Card `json:"cards"`
	PlayTurn int    `json:"playTurn"`
}

type ImportCards struct {
	CardID uint `json:"cardID"`
}
