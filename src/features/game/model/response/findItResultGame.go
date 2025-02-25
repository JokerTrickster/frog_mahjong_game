package response

type ResFindItResult struct {
	Round int          `json:"round"`
	Users []UserResult `json:"users"`
}

type UserResult struct {
	Name              string `json:"name"`
	TotalCorrectCount int    `json:"totalCorrectCount"`
}
