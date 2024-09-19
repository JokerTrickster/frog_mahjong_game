package response

type ResGetUser struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Coin   int    `json:"coin"`
}
