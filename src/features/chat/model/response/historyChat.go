package response

type ResHistoryChat struct {
	Total int           `json:"total"`
	Chats []HistoryChat `json:"chats"`
}

type HistoryChat struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Created string `json:"created"`
	UserID  uint   `json:"userID"`
}