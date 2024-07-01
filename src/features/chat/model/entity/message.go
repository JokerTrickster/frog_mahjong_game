package entity

type MessageEntity struct {
	Secret string `json:"secret"`
}

type AuthEntity struct {
	Secret string `json:"secret"`
	UserID uint   `json:"userID"`
}
