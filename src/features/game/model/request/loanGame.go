package request

type ReqLoan struct {
	LoanUserID int `json:"loanUserID"`
	LoanCardID int `json:"loanCardID"`
	RoomID     int `json:"roomID"`
}
