package request

type ReqReport struct {
	TargetUserID uint   `json:"targetUserID"`
	CategoryID   uint   `json:"categoryID"`
	Reason       string `json:"reason"`
}
