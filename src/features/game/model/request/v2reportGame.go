package request

type ReqV2Report struct {
	TargetUserID uint   `json:"targetUserID"`
	CategoryID   uint   `json:"categoryID"`
	Reason       string `json:"reason"`
}
