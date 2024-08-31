package entity

type WSRequestWinEntity struct {
	RoomID   uint            `json:"roomID omitempty"`
	UserID   uint            `json:"userID"`
	Score    int             `json:"score"`
	Cards    []int           `json:"cards"`
	LoanInfo *ReqWinLoanInfo `json:"loanInfo omitempty"`
}
type ReqWinLoanInfo struct {
	TargetUserID int `json:"targetUserID"`
	CardID       int `json:"cardID"`
}
