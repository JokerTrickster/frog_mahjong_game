package db

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Score    int    `json:"score"`
	State    string `json:"state"`
	RoomID   int    `json:"roomID"`
}
type Rooms struct {
	ID           int    `json:"id"`
	CurrentCount int    `json:"currentCount"`
	MaxCount     int    `json:"maxCount"`
	MinCount     int    `json:"minCount"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	State        string `json:"state"`
	Owner        string `json:"owner"`
}

type RoomUsers struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userID"`
	RoomID      int    `json:"roomID"`
	Score       int    `json:"score"`
	CardCount   int    `json:"cardCount"`
	PlayerState string `json:"playerState"`
}

type Cards struct {
	ID     int    `json:"id"`
	RoomID int    `json:"roomID"`
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	State  string `json:"state"`
}
