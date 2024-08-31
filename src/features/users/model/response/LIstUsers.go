package response

type ResListUser struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

type User struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	State  string `json:"state"`
	Coin   int    `json:"coin"`
}
