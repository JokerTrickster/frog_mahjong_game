package request

type ReqReportImageUploadGame struct {
	FailedList  []string `json:"failedList"`
	SuccessList []string `json:"successList"`
}
