package request

type ReqWSWinEvent struct {
	Cards    []ReqWinCard `json:"cards"`
	PlayTurn int          `json:"playTurn"`
}

type ReqWinCard struct {
	CardID uint   `json:"cardID"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
