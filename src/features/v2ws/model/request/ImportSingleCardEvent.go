package request

type ReqWSImportSingleCard struct {
	CardID   uint `json:"cardID"`
	PlayTurn int  `json:"playTurn"`
}
