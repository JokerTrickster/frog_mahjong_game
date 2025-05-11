package request

type ReqGameOverEvent struct {
	WinnerID uint `json:"winnerID"`
	LoserID  uint `json:"loserID"`
}
