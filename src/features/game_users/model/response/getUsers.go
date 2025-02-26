package response

type ResGetGameUser struct {
	UserID       int    `json:"userID"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Coin         int    `json:"coin"`
	ProfileID    int    `json:"profileID"`
	Disconnected int64    `json:"disconnected"`
}