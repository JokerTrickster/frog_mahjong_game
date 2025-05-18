package request

type ReqWSSuccessEvent struct {
	Cards    []ReqSuccessCard    `json:"cards"`
	Score    int                 `json:"score"`
	LoanInfo *ReqSuccessLoanInfo `json:"loanInfo"`
}

type ReqSuccessCard struct {
	CardID uint `json:"cardID"`
}

type ReqSuccessLoanInfo struct {
	TargetUserID int `json:"targetUserID"`
	CardID       int `json:"cardID"`
}
