package request

type ReqWSWinEvent struct {
	Cards    []ReqWinCard    `json:"cards"`
	Score    int             `json:"score"`
	LoanInfo *ReqWinLoanInfo `json:"loanInfo omitempty"`
}

type ReqWinCard struct {
	CardID uint `json:"cardID"`
}

type ReqWinLoanInfo struct {
	TargetUserID int `json:"targetUserID"`
	CardID       int `json:"cardID"`
}
