package response

type ResFindItResult struct {
	Round int          `json:"round"`
	Users []UserResult `json:"users"`
}

type UserResult struct {
	UserID            int `json:"userID"`
	TotalCorrectCount int `json:"totalCorrectCount"`
}
