package response

type ResV2DrawResult struct {
	Users []DrawResult `json:"users"`
}

type DrawResult struct {
	UserID          int   `json:"userID"`
	SuccessMissions []int `json:"successMissions"`
}
