package entity

type WSSuccessEntity struct {
	RoomID   uint                `json:"roomID omitempty"`
	UserID   uint                `json:"userID"`
	Score    int                 `json:"score"`
	Cards    []int               `json:"cards"`
	LoanInfo *ReqSuccessLoanInfo `json:"loanInfo omitempty"`
}
type ReqSuccessLoanInfo struct {
	TargetUserID int `json:"targetUserID"`
	CardID       int `json:"cardID"`
}
